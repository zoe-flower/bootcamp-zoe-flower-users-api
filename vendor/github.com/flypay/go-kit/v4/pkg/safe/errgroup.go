package safe

import (
	"context"
	"sync"
)

// NOTE: This file is an adapted copy of golang.org/x/sync/errgroup

// A ErrorGroup is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero ErrorGroup is valid and does not cancel on error.
type ErrorGroup struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// ErrorGroupWithContext returns a new ErrorGroup and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func ErrorGroupWithContext(ctx context.Context) (*ErrorGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrorGroup{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *ErrorGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// Unlike (*errgroup.Group).Go, it runs the go routine safely.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *ErrorGroup) Go(f func() error) {
	g.wg.Add(1)

	Go(func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}
