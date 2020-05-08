package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/pkg/errors"
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
		BaseURL:     url,
		HTTPTimeout: DefaultHTTPTimeout,
	}
}

// Client is a simple HTTP client
type Client struct {
	BaseURL     string
	HTTPTimeout time.Duration
}

// GetStateWithOptions gets content for specific key in state store
// TODO: implement with StateOptions
func (c *Client) GetStateWithOptions(ctx context.Context, store, key string, opt *StateOptions) (data []byte, err error) {
	ctx, span := trace.StartSpan(ctx, "get-state")
	defer span.End()

	url := fmt.Sprintf("%s/v1.0/state/%s/%s", c.BaseURL, store, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("consistency", DefaultConsistency)
	req.Header.Set("concurrency", DefaultConcurrency)
	req = req.WithContext(ctx)

	if opt != nil && opt.Concurrency != "" {
		req.Header.Set("concurrency", opt.Concurrency)
	}

	if opt != nil && opt.Consistency != "" {
		req.Header.Set("consistency", opt.Consistency)
	}

	resp, err := c.newHTTPClient().Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error quering state service: %s", url)
	}
	defer resp.Body.Close()

	logger.Printf("%s GET: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	// on initial run there won't be any state
	if resp.StatusCode == http.StatusNoContent ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusUnauthorized {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code from GET to %s: %d", url, resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading response from GET to %s", url)
	}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("store", store),
		trace.StringAttribute("key", key),
	}, "Got state")

	return content, nil

}

// GetState gets content for specific key in state store
func (c *Client) GetState(ctx context.Context, store, key string) (data []byte, err error) {
	return c.GetStateWithOptions(ctx, store, key, nil)
}

// SaveStateData saves state data into state store
func (c *Client) SaveStateData(ctx context.Context, store string, data *StateData) error {
	list := []*StateData{data}
	url := fmt.Sprintf("%s/v1.0/state/%s", c.BaseURL, store)
	return c.post(ctx, "save-state", url, list)
}

// SaveState saves data into state store for specific key
func (c *Client) SaveState(ctx context.Context, store, key string, data interface{}) error {
	state := &StateData{
		Key:   key,
		Value: data,
		Options: &StateOptions{
			Consistency: "strong",     // override default consistency (eventual)
			Concurrency: "last-write", // override defaults (first-write)
		},
		Metadata: map[string]string{
			"created_on": time.Now().UTC().String(),
		},
	}
	return c.SaveStateData(ctx, store, state)
}

// Publish serializes data to JSON and publishes it onto specified topic
func (c *Client) Publish(ctx context.Context, topic string, data interface{}) error {
	url := fmt.Sprintf("%s/v1.0/publish/%s", c.BaseURL, topic)
	return c.post(ctx, "publish-to-topic", url, data)
}

// InvokeBinding serializes data to JSON and submits to specific binding
func (c *Client) InvokeBinding(ctx context.Context, binding string, data interface{}) error {
	url := fmt.Sprintf("%s/v1.0/bindings/%s", c.BaseURL, binding)
	return c.post(ctx, "invoke-service", url, data)
}

// InvokeService serializes input data to JSON and invokes the remote service method
func (c *Client) InvokeService(ctx context.Context, service, method string, in interface{}) (out []byte, err error) {
	ctx, span := trace.StartSpan(ctx, "invoke-service")
	defer span.End()

	url := fmt.Sprintf("%s/v1.0/invoke/%s/method/%s", c.BaseURL, service, method)
	b, _ := json.Marshal(in)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := c.newHTTPClient().Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error invoking service: %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code from GET to %s: %d", url, resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading response from invoke to %s", url)
	}

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("service", service),
		trace.StringAttribute("method", method),
	}, "Invoked service")

	return content, nil
}

func (c *Client) post(ctx context.Context, method, url string, data interface{}) error {
	ctx, span := trace.StartSpan(ctx, method)
	defer span.End()

	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := c.newHTTPClient().Do(req)
	if err != nil {
		return errors.Wrapf(err, "error posting %+v to %s", data, url)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with result: %+v - %q",
			url, data, dump)
	}

	return nil
}

func (c *Client) newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.HTTPTimeout,
	}
}
