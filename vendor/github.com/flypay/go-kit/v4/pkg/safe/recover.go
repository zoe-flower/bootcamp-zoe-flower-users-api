package safe

import (
	"runtime/debug"
	"time"

	"github.com/flypay/go-kit/v4/pkg/log"
	"github.com/flypay/go-kit/v4/pkg/runtime"
	"github.com/getsentry/sentry-go"
)

// Recover will recover from any panics will either re-raise the panic
// if configured to do so, otherwise it will log the recover and send
// to sentry
func Recover() {
	if r := recover(); r != nil {
		if cfg := runtime.DefaultConfig(); cfg.PanicOnRecover {
			panic(r)
		} else {
			log.Errorf("handling recover for: %+v:\n%s", r, debug.Stack())
			hub := sentry.CurrentHub()
			hub.Recover(r)
			hub.Flush(time.Second * 5)
		}
	}
}
