package qvapay

import (
	"context"
	"io"
	"net/http"
	"os"
)

type apiClient struct {
	// endpoint http
	server string
	client *http.Client
	//optional for debuging
	debug io.Writer
}

// HttpRequestDoer defines methods for a http client.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// APIClientOptions contains optios used to initialize an API Client using
// NewAPIClient
type APIClientOptions struct {
	//optional, defaults to http.DefaultClient
	HttpClient *http.Client
	//optional for debuging
	Debug io.Writer
	// Server url to use
	Server string
}

// NewAPICLient creates a client from APIClientOptions. If no options are provided,
//  client with default configurations is returned.
func NewAPIClient(options APIClientOptions) APIClient {
	c := &apiClient{}
	if options.Server == "" {
		options.Server = os.Getenv("QVAPAY_API_ENDPOINT")
	}

	c.server = options.Server
	c.client = options.HttpClient
	c.debug = options.Debug
	return c

}

// APIClient defines methods implemented by the api client.
type APIClient interface {
	Login(ctx context.Context, payload LoginRequest) (APIResult, error)
}
