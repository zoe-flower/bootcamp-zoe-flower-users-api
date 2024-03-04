package safe

// Go is a safe way to run goroutines. It will recover from a runtime panic and
// it will log and send the resulting panic to Sentry.
func Go(f func()) {
	go func() {
		// Recover and send to Sentry
		defer Recover()
		// Invoke the function
		f()
	}()
}
