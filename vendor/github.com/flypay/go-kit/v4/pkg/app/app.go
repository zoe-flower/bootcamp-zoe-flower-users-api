package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/flypay/go-kit/v4/pkg/env"
	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/safe"
)

// Period after which services are shutdown
const kubernetesGracePeriod = 15 * time.Second

// Amount of time to allow a service to gracefully shutdown
const shutdownTimeout = 6 * time.Second

var closeables []Closeable

// GracefullyShutdown adds a type that we must wait for to close before the application is terminated.
// This is commonly used for the eventbus producer and consumer so that in-flight messages are handled before termination.
func GracefullyShutdown(services ...Closeable) {
	closeables = append(closeables, services...)
}

// Closeable is a system that needs to be closed gracefully when the
// application is signalled to terminate
type Closeable interface {
	Close() error
}

// CloseableFunc can be used to inline a closeable that doesn't implement the close interface
type CloseableFunc func() error

// Close calls the CloseableFunc
func (c CloseableFunc) Close() error {
	return c()
}

// Shutdowner is a capability to gracefully shutdown a service. A context is required so that we
// can give it a timeout after which anything still waiting will be forcefully terminated.
type Shutdowner interface {
	Shutdown(context.Context) error
}

// Run will start the application and wait for termination signals.
// Any services passed in with be closed (or shutdown gracefully see Shutdowner) when the application
// receives a termination signal
func Run() {
	log.Info("Running app")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Infof("App received signal: %s", <-c)

	// It takes several seconds to remove the service from the http traffic routing
	// so we wait before shutting down
	if !env.IsLocal() {
		time.Sleep(kubernetesGracePeriod)
	}

	log.Infof("Attempting to stop services")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(closeables))

	for _, svc := range closeables {
		func(svc Closeable) {
			safe.Go(func() {
				shutdown(ctx, wg, svc)
			})
		}(svc)
	}

	log.Infof("App waiting for processes to close")
	wg.Wait()

	log.Infof("App closed")
}

func logErr(err error, msg string) {
	if err != nil {
		log.Errorf(msg, err)
	}
}

func shutdown(ctx context.Context, wg *sync.WaitGroup, svc Closeable) {
	defer wg.Done()

	// Try graceful shutdown if the service implements it
	if shutter, ok := svc.(Shutdowner); ok {
		logErr(shutter.Shutdown(ctx), "Error shutting down service: %+v")
		return
	}

	logErr(svc.Close(), "Error stopping service: %+v")
}
