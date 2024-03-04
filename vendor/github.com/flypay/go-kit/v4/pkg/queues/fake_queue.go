package queues

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// fakeJob can be used for testing things that use a job
type fakeJob struct {
	identifier int
	deleted    bool
}

func (j *fakeJob) ID() string {
	return fmt.Sprintf("%d", j.identifier)
}

func (j *fakeJob) Payload() io.Reader {
	return &bytes.Buffer{}
}

func (j *fakeJob) Delete() {
	j.deleted = true
}

func (j *fakeJob) Received() time.Time {
	return time.Now()
}

func (j *fakeJob) Context() context.Context {
	return context.Background()
}

func (j *fakeJob) Attributes() map[string]*sqs.MessageAttributeValue {
	return make(map[string]*sqs.MessageAttributeValue)
}

func (j *fakeJob) RawMessage() []byte { return nil }

// fakeQueue can be used for testing things that use a queue
type fakeQueue struct {
	jobs         []Job
	delayedJobs  chan Job
	receiveErr   error
	recieveCount int64
	lock         sync.RWMutex
}

func NewFakeQueue() Queue { return &fakeQueue{} }

func (mq *fakeQueue) ReceiveCount() int64 {
	mq.lock.RLock()
	defer mq.lock.RUnlock()
	return mq.recieveCount
}

func (mq *fakeQueue) Delay(j Job, duration time.Duration) error {
	go func() {
		mq.delayedJobs <- j
	}()
	return nil
}

func (mq *fakeQueue) Delete(j Job) error {
	return nil
}

func (mq *fakeQueue) SeedJobs(jobs []Job) {
	mq.jobs = jobs
}

func (mq *fakeQueue) ReceiveWillErrorWith(err error) {
	mq.receiveErr = err
}

func (mq *fakeQueue) Receive(max int) ([]Job, error) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.recieveCount++

	if len(mq.jobs) == 0 {
		return []Job{}, mq.receiveErr
	}

	jobs := mq.jobs[:max]
	mq.jobs = mq.jobs[max:]

	return jobs, nil
}

func (mq *fakeQueue) Stop() {}

func (mq *fakeQueue) SendMessageBatch(events []string) error {
	mq.lock.Lock()
	defer mq.lock.Unlock()

	for i := range events {
		mq.jobs = append(mq.jobs, &fakeJob{i, false})
	}

	return nil
}

func (mq *fakeQueue) SendMessage(body string) error {
	mq.lock.Lock()
	defer mq.lock.Unlock()

	mq.jobs = append(mq.jobs, &fakeJob{0, false})

	return nil
}
