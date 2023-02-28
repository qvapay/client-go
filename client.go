package qvapay

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
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
//
//	client with default configurations is returned.
func NewAPIClient(options APIClientOptions) APIClient {
	c := &apiClient{}
	if options.Server == "" {
		options.Server = os.Getenv("QVAPAY_API_ENDPOINT")
	}

	c.server = options.Server
	if options.HttpClient != nil {
		c.client = options.HttpClient
	} else {
		c.client = http.DefaultClient
	}

	// Logging consideration
	c.client.Transport = &authRoundTripper{
		next: &headersRoundTripper{
			next: &debugRoundTripper{
				next:   http.DefaultTransport,
				logger: os.Stdout,
				Debug:  options.Debug,
			},
		},
	}

	return c

}

// RoundTrippers Section
type authRoundTripper struct {
	next http.RoundTripper
}

func (t authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	switch r.URL.Path {
	case "/api/auth/login", "/api/auth/register":
		return t.next.RoundTrip(r)
	default:
		if authUser == nil {
			return t.next.RoundTrip(r)
		}

		if authUser.AccessToken != "" {
			r.Header.Set("Authorization", "Bearer "+authUser.AccessToken)
			return t.next.RoundTrip(r)

		} else {
			return t.next.RoundTrip(r)
		}
	}

}

type debugRoundTripper struct {
	next   http.RoundTripper
	logger io.Writer
	Debug  io.Writer
}

func (l *debugRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Here we can set a fancy Logging implementation
	if l.Debug != nil {
		// Request Dump
		dump, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			fmt.Printf("error dumping HTTP request: %s", err.Error())
			return l.next.RoundTrip(r)
		}
		fmt.Fprintln(l.Debug, string(dump))
		fmt.Fprintln(l.Debug)

		// Response Dump
		resp, err := l.next.RoundTrip(r)
		dump, _ = httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Printf("error dumping HTTP reponse: %s", err.Error())
			return l.next.RoundTrip(r)
		}
		fmt.Fprintln(l.Debug, string(dump))
		fmt.Fprintln(l.Debug)

		return resp, nil
	}

	return l.next.RoundTrip(r)
}

type headersRoundTripper struct {
	next http.RoundTripper
}

func (l *headersRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("content-type", "application/json")
	r.Header.Add("User-Agent", "qvapay-go")

	return l.next.RoundTrip(r)
}

func DrainBody(respBody io.ReadCloser) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if respBody != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		defer func(respBody io.ReadCloser) {
			err := respBody.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(respBody)
		_, _ = io.Copy(ioutil.Discard, respBody)
	}
}

// APIClient defines methods implemented by the api client.
type APIClient interface {
	Login(ctx context.Context, payload LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, payload RegisterRequest) (RegisterResponse, error)
	Logout(ctx context.Context) (LogoutResponse, error)
	GetMeRAW(ctx context.Context) (MeRAW, error)
	GetMe(ctx context.Context) (User, error)
}
