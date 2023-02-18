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
	c.client = options.HttpClient

	// Logging consideration
	c.client.Transport = &headersRoundTripper{
		next: &loggingRoundTripper{
			next:   http.DefaultTransport,
			logger: os.Stdout,
		},
	}

	c.debug = options.Debug
	return c

}

// APIClient defines methods implemented by the api client.
type APIClient interface {
	Login(ctx context.Context, payload LoginRequest) (APIResult, error)
}

// RoundTrippers Section
type authRoundTripper struct {
	next  http.RoundTripper
	Token string
}

func (t authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.Token)
	return t.next.RoundTrip(r)
}

type loggingRoundTripper struct {
	next   http.RoundTripper
	logger io.Writer
}

func (l *loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Here we can set a fancy Logging implementation
	// Request Dump
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		fmt.Printf("error dumping HTTP request: %s", err.Error())
		return l.next.RoundTrip(r)
	}
	log.Println("Request:", string(dump))

	// Response Dump
	resp, err := l.next.RoundTrip(r)
	dump, _ = httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf("error dumping HTTP reponse: %s", err.Error())
		return l.next.RoundTrip(r)
	}
	log.Println("Response:", string(dump))

	return resp, nil
}

type headersRoundTripper struct {
	next http.RoundTripper
}

func (l *headersRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("content-type", "application/json")
	r.Header.Add("User-Agent", "qvapay-go")

	return l.next.RoundTrip(r)
}

// // dumpResponse writes the raw response data to the debug output, if set, or
// // standard error otherwise.
// func (c *apiClient) dumpResponse(resp *http.Response) {
// 	// ignore errors dumping response - no recovery from this
// 	responseDump, err := httputil.DumpResponse(resp, true)
// 	if err != nil {
// 		log.Fatalf("dumpResponse: " + err.Error())
// 	}
// 	fmt.Fprintln(c.debug, string(responseDump))
// 	fmt.Fprintln(c.debug)
// }

// func apiCallDebugger(req *http.Request, debug io.Writer) error {
// 	req.Header.Add("content-type", "application/json")
// 	req.Header.Add("User-Agent", "qvapay-go")
// 	if debug != nil {
// 		requestDump, err := httputil.DumpRequestOut(req, true)
// 		if err != nil {
// 			return fmt.Errorf("error dumping HTTP request: %v", err)
// 		}
// 		fmt.Fprintln(debug, string(requestDump))
// 		fmt.Fprintln(debug)
// 	}
// 	return nil
// }

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
