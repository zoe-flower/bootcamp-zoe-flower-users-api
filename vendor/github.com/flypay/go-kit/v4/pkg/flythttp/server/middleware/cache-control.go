package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Background: According to JET's API Principles and Practices, a server
// "MUST provide a Cache-Control header to indicate to clients and intermediary servers if and how long they can cache the data for"
// See: https://pages.github.je-labs.com/Architecture/principles-and-practices/practices/http-api-design/api-response-caching/

const cacheControl = "Cache-Control"

// values taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	// The following is taken from https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control:
	// - "The no-cache response directive indicates that the response can be stored in caches, but must be validated with the origin server before each reuse — even when the cache is disconnected from the origin server."
	// - "The private response directive indicates that the response can be stored only in a private cache (e.g. local caches in browsers).""
	// - "max-age=0 is a workaround for no-cache, because many old (HTTP/1.0) cache implementations don't support no-cache."
	cacheControl: "no-cache, private, max-age=0",

	// "The Expires HTTP header contains the date/time after which the response is considered expired.
	// Invalid expiration dates with value 0 represent a date in the past and mean that the resource is already expired."
	"Expires": time.Unix(0, 0).UTC().Format(time.RFC1123),

	// "The Pragma HTTP/1.0 general header is an implementation-specific header that may have various effects along the request-response chain. This header serves for backwards compatibility with the HTTP/1.0 caches that do not have a Cache-Control HTTP/1.1 header."
	// If set to 'no-cache':
	// "Same as Cache-Control: no-cache. Forces caches to submit the request to the origin server for validation before a cached copy is released"
	"Pragma": "no-cache",
}

// CacheControl adds a Cache-Control header (no cache) if not already set.
func CacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		err = next(c)

		h := c.Response().Header()

		if h.Get(cacheControl) != "" {
			// is already set, don't overwrite
			return
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			h.Set(k, v)
		}

		return
	}
}
