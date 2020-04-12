package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/pkg/errors"
)

var (
	logger             = log.New(os.Stdout, "CLIENT == ", 0)
	defaultHTTPTimeout = time.Second * 30
	defaultConsistency = "strong"     // override defaults (eventual)
	defaultConcurrency = "last-write" // override defaults (first-write)
)

// NewClient creates valid instance of Client
// baseURL (e.g. http://localhost:3500)
// API version and function (state, publish) will be added by client
func NewClient(baseURL string) (client *Client) {
	return &Client{
		BaseURL:     baseURL,
		HTTPTimeout: defaultHTTPTimeout,
	}
}

// Client is a simple HTTP client
type Client struct {
	BaseURL     string
	HTTPTimeout time.Duration
}

// GetDataWithOptions gets content for specific key in state store
// TODO: implement with StateOptions
func (c *Client) GetDataWithOptions(store, key string, opt *StateOptions) (data []byte, err error) {

	url := fmt.Sprintf("%s/v1.0/state/%s/%s", c.BaseURL, store, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("consistency", defaultConsistency)
	req.Header.Set("concurrency", defaultConcurrency)

	if opt != nil && opt.Concurrency != "" {
		req.Header.Set("concurrency", opt.Concurrency)
	}

	if opt != nil && opt.Consistency != "" {
		req.Header.Set("consistency", opt.Consistency)
	}

	//TODO: Handle RetryPolicy

	resp, err := c.newHTTPClient().Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error quering state service: %s", url)
	}
	defer resp.Body.Close()

	logger.Printf("%s GET: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	// on initial run there won't be any state
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code from GET to %s: %d", url, resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading response from GET to %s", url)
	}

	return content, nil

}

// GetData gets content for specific key in state store
func (c *Client) GetData(store, key string) (data []byte, err error) {
	return c.GetDataWithOptions(store, key, nil)
}

// Save saves state data into state store
// TODO: check result of publish and consider returning
func (c *Client) Save(store string, data *StateData) error {

	b, _ := json.Marshal([]*StateData{data})

	url := fmt.Sprintf("%s/v1.0/state/%s", c.BaseURL, store)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.newHTTPClient().Do(req)
	if err != nil {
		return errors.Wrapf(err, "error posting to %s with key: %s, data: %v", url, data.Key, data)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusCreated {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with key: %s, data: %v - %q",
			url, data.Key, data, dump)
	}

	return nil

}

// SaveData saves data into state store for specific key
func (c *Client) SaveData(store, key string, data interface{}) error {

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

	return c.Save(store, state)

}

// Publish serializes data to JSON and publishes it onto specified topic
// TODO: check result of publish and consider returning
func (c *Client) Publish(topic string, data interface{}) error {

	url := fmt.Sprintf("%s/v1.0/publish/%s", c.BaseURL, topic)

	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.newHTTPClient().Do(req)
	if err != nil {
		return errors.Wrapf(err, "error publishing result %+v to %s", data, url)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
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
