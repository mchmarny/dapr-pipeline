package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

var (
	logger = log.New(os.Stdout, "CLIENT == ", 0)

	// DefaultHTTPTimeout is the default HTTP client timeout
	DefaultHTTPTimeout = time.Second * 30

	// DefaultConsistency is the state store consistency option setting
	DefaultConsistency = "eventual" // override defaults (eventual)

	// DefaultConcurrency is the state store concurrency option setting
	DefaultConcurrency = "last-write" // override defaults (first-write)

	// DefaultRetryPolicyInterval is the state store retry policy interval setting
	DefaultRetryPolicyInterval = 100

	// DefaultRetryPolicyThreshold is the state store retry policy threshold setting
	DefaultRetryPolicyThreshold = 3

	// DefaultRetryPolicyPattern is the state store retry policy pattern setting
	DefaultRetryPolicyPattern = "exponential"
)

// NewClient creates instance of dapr Client using http://localhost:PORT
// where PORT is the value of DAPR_HTTP_PORT env var defaulted to 3500
func NewClient() (client *Client) {
	port := os.Getenv("DAPR_HTTP_PORT")
	if port == "" {
		port = "3500"
	}
	url := fmt.Sprintf("http://localhost:%s", port)
	return NewClientWithURL(url)
}

// NewClientWithURL creates valid instance of Client using provided url
func NewClientWithURL(url string) (client *Client) {
	if url == "" {
		url = "http://localhost:3500"
	}
	return &Client{
		url:     url,
		timeout: DefaultHTTPTimeout,
	}
}

// Client is a simple HTTP Dapr client abstraction
type Client struct {
	url     string
	timeout time.Duration
}

func (c *Client) exec(ctx trace.SpanContext, req *http.Request) (out []byte, status int, err error) {
	if req == nil {
		err = errors.New("nil request")
		return
	}

	httpFmt := tracecontext.HTTPFormat{}
	httpFmt.SpanContextToRequest(ctx, req)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout:   c.timeout,
		Transport: &ochttp.Transport{},
	}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "error executing %+v", req)
		return
	}
	defer resp.Body.Close()

	status = resp.StatusCode
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "error reading response from %+v", resp)
		return
	}
	out = content

	return
}
