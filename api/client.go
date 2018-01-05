package api

import "net/http"

// Client is an interface that represents the API client of Loggly.
type Client interface {
	Log(body []byte) (*http.Response, error)
	LogAsBulk(body []byte) (*http.Response, error)
	SetHTTPClient(client *http.Client)
}
