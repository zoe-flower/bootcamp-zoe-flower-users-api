package queues

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	flytcontext "github.com/flypay/go-kit/v4/pkg/context"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/google/uuid"
)

const (
	// flytRequestID is the key against which the request id is stored in message attributes
	flytRequestID string = "X-Flyt-Request-ID"

	// maxReceiveMessages is the limit of messages SQS will return in one Receive request.
	maxReceiveMessages int = 10

	// allAttributes informs SQS to return all message attributes
	allAttributes = "All"
)

type sqsmq struct {
	queueURL          string
	sqs               sqsiface.SQSAPI
	waitTimeSeconds   int64
	visibilityTimeOut time.Duration

	stop chan struct{}
}

// NewSQS returns a new Queue implementation for Amazon SQS
func NewSQS(url string, sess *session.Session, waitTimeSeconds int64, visibilityTimeOut time.Duration) Queue {
	s := sqs.New(sess)
	s.AddDebugHandlers()

	q := sqsmq{
		queueURL:          url,
		sqs:               s,
		waitTimeSeconds:   waitTimeSeconds,
		visibilityTimeOut: visibilityTimeOut,
		stop:              make(chan struct{}),
	}

	return &q
}

// NewQueue accepts an SQS interface to inject into it
func NewQueue(url string, sqs sqsiface.SQSAPI, waitTimeSeconds int64, visibilityTimeOut time.Duration) Queue {
	return &sqsmq{
		queueURL:          url,
		sqs:               sqs,
		waitTimeSeconds:   waitTimeSeconds,
		visibilityTimeOut: visibilityTimeOut,
		stop:              make(chan struct{}),
	}
}

func (q *sqsmq) Receive(max int) ([]Job, error) {
	var (
		err error
		res *sqs.ReceiveMessageOutput
	)

	if max > maxReceiveMessages {
		max = maxReceiveMessages
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Ensure we always cancel after completing the Receive
	defer cancel()

	// If Queue.Stop() is called whilst waiting for messages we cancel the current context
	go func() {
		select {
		case <-q.stop:
			cancel()
		case <-ctx.Done():
			// Ensure we do not leak goroutines when the Receive context completes
			return
		}
	}()

	res, err = q.sqs.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              &q.queueURL,
		WaitTimeSeconds:       aws.Int64(q.waitTimeSeconds),
		MaxNumberOfMessages:   aws.Int64(int64(max)),
		MessageAttributeNames: []*string{aws.String(allAttributes)},
		VisibilityTimeout:     aws.Int64(int64(q.visibilityTimeOut.Seconds())),
	})
	if err != nil {
		// When stopped no error should be returned as we are exiting.
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, nil
		}
		return nil, err
	}

	jobs := make([]Job, 0, len(res.Messages))
	for _, msg := range res.Messages {
		j := newSQSJob(context.Background(), msg, q, q.visibilityTimeOut)

		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (q *sqsmq) Stop() {
	close(q.stop)
}

func (q *sqsmq) Delete(j Job) error {
	job, ok := j.(*sqsjob)
	if !ok {
		return fmt.Errorf("job was not an implementation of *sqsjob: %s: %s", reflect.TypeOf(j), j.ID())
	}

	_, err := q.sqs.DeleteMessageWithContext(job.ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueURL,
		ReceiptHandle: job.ReceiptHandle,
	})

	return err
}

func (q *sqsmq) Delay(j Job, visibilityTimeout time.Duration) error {
	job, ok := j.(*sqsjob)
	if !ok {
		return fmt.Errorf("job was not an implementation of *sqsjob: %s: %s", reflect.TypeOf(j), j.ID())
	}

	timeoutSeconds := int64(visibilityTimeout.Seconds())

	_, err := q.sqs.ChangeMessageVisibilityWithContext(
		job.ctx,
		&sqs.ChangeMessageVisibilityInput{
			QueueUrl:          &q.queueURL,
			ReceiptHandle:     job.ReceiptHandle,
			VisibilityTimeout: &timeoutSeconds,
		})
	if err != nil {
		return err
	}

	return nil
}

type SendMessageBatchError struct {
	Failed []*sqs.BatchResultErrorEntry
	Err    error
}

func (s *SendMessageBatchError) Error() string {
	if len(s.Failed) > 0 {
		return fmt.Sprintf("Failed to send %d messages", len(s.Failed))
	}
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

// SendMessageBatch will send only first 10 messages out of the 'messages' argument
// This is a limitation of the aws sqs api, so you shouldn't be passing more than
// 10 messages to this method.
func (q *sqsmq) SendMessageBatch(messages []string) error {
	numEntries := 10
	if l := len(messages); l < numEntries {
		numEntries = l
	}

	entries := make([]*sqs.SendMessageBatchRequestEntry, 0, numEntries)
	for i := 0; i < numEntries; i++ {
		entries = append(entries, &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(uuid.NewString()),
			MessageBody: aws.String(messages[i]),
		})
	}

	params := &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(q.queueURL),
	}
	output, err := q.sqs.SendMessageBatch(params)
	if err != nil {
		return &SendMessageBatchError{Err: err}
	}
	if len(output.Failed) > 0 {
		return &SendMessageBatchError{Failed: output.Failed}
	}

	return nil
}

func (q *sqsmq) SendMessage(body string) error {
	msg := &sqs.SendMessageInput{
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(q.queueURL),
	}

	_, err := q.sqs.SendMessage(msg)
	return err
}

type sqsjob struct {
	*sqs.Message

	queue             Queue
	received          time.Time
	ctx               context.Context
	visibilityTimeout time.Duration
}

func newSQSJob(ctx context.Context, m *sqs.Message, mq Queue, visibilityTimeOut time.Duration) *sqsjob {
	if rid, ok := m.MessageAttributes[flytRequestID]; ok {
		ctx = flytcontext.WithRequestID(ctx, *rid.StringValue)
	} else {
		// Fallback to using the message identifier
		ctx = flytcontext.WithRequestID(ctx, *m.MessageId)
	}

	job := &sqsjob{
		Message: m,
		queue:   mq,
		// TODO received time should come from message attributes
		received:          time.Now(),
		ctx:               ctx,
		visibilityTimeout: visibilityTimeOut,
	}

	return job
}

func (j *sqsjob) ID() string {
	return *j.MessageId
}

func (j *sqsjob) Attributes() map[string]*sqs.MessageAttributeValue {
	return j.MessageAttributes
}

func (j *sqsjob) Payload() io.Reader {
	return strings.NewReader(*j.Body)
}

func (j *sqsjob) Delete() {
	logger := log.WithContext(j.Context())

	if err := j.queue.Delete(j); err != nil {
		logger.Errorf("Failed to delete job: %s, error: %s", j.ID(), err)
		return
	}
}

func (j *sqsjob) Received() time.Time {
	return j.received
}

func (j *sqsjob) Context() context.Context {
	return j.ctx
}

func (j *sqsjob) RawMessage() []byte {
	b, err := json.Marshal(j.Message)
	if err != nil {
		log.WithContext(j.Context()).Errorf("Failed to marshal job: %s, error: %s", j.ID(), err)
	}
	return b
}
