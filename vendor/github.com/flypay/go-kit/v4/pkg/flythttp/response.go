package flythttp

import "net/http"

func IsServerError(response *http.Response) bool {
	return response.StatusCode >= http.StatusInternalServerError && response.StatusCode <= http.StatusNetworkAuthenticationRequired
}

func IsClientError(response *http.Response) bool {
	return response.StatusCode >= http.StatusBadRequest && response.StatusCode <= http.StatusUnavailableForLegalReasons
}
