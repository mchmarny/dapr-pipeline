package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/pkg/errors"
)

// SimpleStateData represents simplified dapr state item
type SimpleStateData struct {
	Key     string              `json:"key"`
	Value   interface{}         `json:"value"`
	Options *SimpleStateOptions `json:"options,omitempty"`
}

// SimpleStateOptions is the dapr state data option for StateData
type SimpleStateOptions struct {
	Consistency string `json:"consistency,omitempty"`
}

// NewClient creates instance of Client
func NewClient(baseURL string) (client *Client) {
	return &Client{
		BaseURL:                baseURL,
		StrongStateConsistency: true,
		Timeout:                time.Second * 30,
	}
}

// Client is a simple HTTP client
type Client struct {
	BaseURL                string
	StrongStateConsistency bool
	Timeout                time.Duration
}

// GetState gets content for specific key in state store
func (c *Client) GetState(store, key string) (data []byte, err error) {

	url := fmt.Sprintf("%s/v1.0/state/%s/%s", c.baseURL, store, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if c.StrongStateConsistency {
		req.Header.Set("consistency", "strong") // TODO: parameterize
	}

	resp, err := newHTTPClient().Do(req)
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

// SaveState saves data into state store for specific key
func (c *Client) SaveState(store, key string, data interface{}) error {

	state := &SimpleStateData{
		Key:     key,
		Value:   data,
		Options: &SimpleStateOptions{Consistency: "strong"},
	}
	b, _ := json.Marshal([]*SimpleStateData{state})

	url := fmt.Sprintf("%s/v1.0/state/%s", c.baseURL, store)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newHTTPClient().Do(req)
	if err != nil {
		return errors.Wrapf(err, "error posting to %s with key: %s, data: %v", url, key, data)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusCreated {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with key: %s, data: %v - %q",
			stateURL, key, data, dump)
	}

	return nil

}

// Publish serializes data to JSON and publishes it onto specified topic
func (c *Client) Publish(topic string, data interface{}) error {

	url := fmt.Sprintf("%s/v1.0/publish/%s", c.baseURL, topic)

	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newHTTPClient().Do(req)
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

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
	}
}
