package queues

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/safe"
	"github.com/go-kit/kit/metrics"
)

// Queue is the interface for message queues
type Queue interface {
	QueueReceiver
	Delete(Job) error
	Delay(Job, time.Duration) error
	SendMessageBatch([]string) error
	SendMessage(string) error
}

// QueueReceiver exposes a function to receive jobs
type QueueReceiver interface {
	Receive(max int) ([]Job, error)
	Stop()
}

type ReceiverFunc func(max int) ([]Job, error)

func (f ReceiverFunc) Receive(max int) ([]Job, error) {
	return f(max)
}

func (f ReceiverFunc) Stop() {
	// Unable to stop arbitrary funcs
}

// Job is a unit of work such as a message received on a queue
type Job interface {
	ID() string
	Payload() io.Reader
	Delete()
	Received() time.Time
	Context() context.Context
	Attributes() map[string]*sqs.MessageAttributeValue
	RawMessage() []byte
}

// Producer is responsible for coordinating the application process
type Producer struct {
	Pool      chan chan Job
	queue     QueueReceiver
	stop      chan struct{}
	done      chan struct{}
	errorWait time.Duration
}

// NewProducer creates and starts a producer which will start filling the job channel with jobs
// A process is started in the producer that continuously receives jobs from the queue
// The producer contains a worker pool with the given capacity, the same number of Consumers should be started.
func NewProducer(q QueueReceiver, poolSize int, capacity metrics.Gauge) *Producer {
	p := &Producer{
		Pool:      make(chan chan Job, poolSize),
		queue:     q,
		stop:      make(chan struct{}),
		done:      make(chan struct{}),
		errorWait: 1 * time.Second,
	}

	safe.Go(func() {
		p.monitor(capacity)
	})
	safe.Go(p.run)

	return p
}

func (p *Producer) run() {
	defer close(p.done)

	for {
		// Always try and read from the stop channel first. This will prioritise
		// stopping over receiving workers off the pool.
		select {
		case <-p.stop:
			return
		default:
			// go and try the pool
		}

		select {
		// A worker is available in the pool
		case worker := <-p.Pool:
			// Receive as many messages as the number of workers waiting
			// in the pool plus the worker we already hold
			jobs, err := p.queue.Receive(len(p.Pool) + 1)
			if err != nil {
				p.handleErr(err)
				p.Pool <- worker
				break
			}
			p.resetErr()

			if len(jobs) == 0 {
				// No work to do, put the worker back on the pool
				p.Pool <- worker
				break
			}

			for i, job := range jobs {
				worker <- job
				// if there are more jobs get a new worker from the pool
				if i < len(jobs)-1 {
					worker = <-p.Pool
				}
			}
		case <-p.stop:
			return
		}
	}
}

func (p *Producer) monitor(capcity metrics.Gauge) {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			capcity.Set(float64(len(p.Pool)))
		case <-p.stop:
			return
		}
	}
}

// Stop gracefull stops the dispatcher process
func (p *Producer) Stop() {
	close(p.stop)
	p.queue.Stop()
	<-p.done
}

// handleErr will block for a period of time if we're experiencing errors receiving from the queue
// when the block time exceeds the maximum the application will restart -- not ideal
func (p *Producer) handleErr(err error) {
	if p.errorWait > 1*time.Minute {
		log.Fatalf("Producer has been erroring for %s, the service will be stopped: %s", p.errorWait, err)
		// Fatal above will terminate the application immediately, posibility that inflight messages will be lost
	}

	log.Warnf("Producer is sleeping for %s as it encountered an error reading from the queue: %s", p.errorWait, err)
	select {
	case <-time.After(p.errorWait):
		p.errorWait *= 2
		return
	case <-p.stop:
		return
	}
}

func (p *Producer) resetErr() {
	p.errorWait = 1 * time.Second
}

// InstrumentedQueueReceiver Wraps a queue with a rate metric for measuring receive rate
type InstrumentedQueueReceiver struct {
	QueueReceiver
	// TODO deprecate the rate as a histogram also provides this metric
	Rate     metrics.Counter
	Quantity metrics.Histogram
}

// Receive gets jobs from the queue and records success and error rate and number of jobs received
func (q *InstrumentedQueueReceiver) Receive(max int) ([]Job, error) {
	jobs, err := q.QueueReceiver.Receive(max)
	state := errorOrSuccess(err)
	q.Rate.With("status", state).Add(1)
	q.Quantity.With("status", state).Observe(float64(len(jobs)))
	return jobs, err
}

func errorOrSuccess(err error) string {
	if err != nil {
		return "error"
	}
	return "success"
}
