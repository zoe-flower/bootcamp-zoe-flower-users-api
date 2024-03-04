package flythttp

import (
	"fmt"

	"github.com/flypay/go-kit/v4/pkg/runtime"
)

// ResolveServiceHost will return the services hostname depending on
// the environment, eg:
// - production: https://<serviceName>.flyt-platform.com
// - staging: https://<serviceName>.flyt-staging.com
func ResolveServiceHost(serviceName string) string {
	cfg := runtime.DefaultConfig()
	// If we have a service hostname override then we should return that
	// without any string formatting.
	// This is mostly used when being run locally.
	if cfg.ServiceHostnameOverride != "" {
		return cfg.ServiceHostnameOverride
	}
	return fmt.Sprintf(cfg.ServiceHostnameFormat, serviceName)
}
