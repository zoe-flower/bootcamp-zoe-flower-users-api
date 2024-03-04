package queues

import (
	"sync"

	"github.com/flypay/go-kit/v4/pkg/log"

	"github.com/flypay/go-kit/v4/pkg/safe"
)

// Consumer pulls work off the jobs channel and performs work on them
type Consumer struct {
	pool   chan chan Job
	jobs   chan Job
	handle Handler
	done   chan struct{}
}

// Handler is the function that processes jobs
// It is expected that the handler should call `Delete()` or `Stop()` on a job after successfully handling it.
// When the returned bool is true the handler was successful, false indicates that the handler failed but was not an error scenario
// Eventually the bool should be removed in favour of error types.
type Handler func(Job) (bool, error)

// Middleware is a function for wrapping a handler with metrics and logging concerns
type Middleware func(Handler) Handler

// NewConsumer creates a running consumer
// The consumer registers itself with the producer pool and executes the handler when receiving a job
func NewConsumer(pool chan chan Job, handler Handler) *Consumer {
	c := &Consumer{
		pool:   pool,
		jobs:   make(chan Job),
		handle: handler,
		done:   make(chan struct{}),
	}

	safe.Go(func() {
		c.run()
	})

	return c
}

// Stop gracefully waits for inflight work to complete before returning
func (c *Consumer) Stop() {
	c.done <- struct{}{}
}

func (c *Consumer) run() {
	for {
		// consumer is ready for work, place on queue
		c.pool <- c.jobs
		select {
		// wait for some work to come in
		case j := <-c.jobs:
			_, err := c.handle(j)
			if err != nil {
				log.WithContext(j.Context()).Errorf("Failed to handle job: %s, error: %s", j.ID(), err)
			}
		case <-c.done:
			return
		}
	}
}

// Stoppable is an interface for something that can be stopped
type Stoppable interface {
	Stop()
}

// StopAll calls Stop() on all the given stoppables and blocks until they all return
func StopAll(stoppables ...Stoppable) {
	wg := sync.WaitGroup{}

	wg.Add(len(stoppables))

	for _, s := range stoppables {
		func(s Stoppable) {
			safe.Go(func() {
				defer wg.Done()
				s.Stop()
			})
		}(s)
	}

	wg.Wait()
}
