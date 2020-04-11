package main

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

var (
	clientTimeout = time.Second * 30
)

// StateData represents dapr state item
type StateData struct {
	Key      string            `json:"key"`
	Value    string            `json:"value"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Options  *StateOptions     `json:"options,omitempty"`
}

// StateOptions is the dapr state data options
type StateOptions struct {
	Concurrency string       `json:"concurrency,omitempty"`
	Consistency string       `json:"consistency,omitempty"`
	RetryPolicy *RetryPolicy `json:"retryPolicy,omitempty"`
}

// RetryPolicy is the dapr StateOptions retry policy
type RetryPolicy struct {
	Threshold int32  `json:"threshold,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	// Interval  *duration.Duration `json:"interval,omitempty"`
}

func getLastID(key string) (id string, err error) {

	url := fmt.Sprintf("%s/%s", stateURL, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("consistency", "strong")

	client := &http.Client{
		Timeout: clientTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "error quering state service: %s", url)
	}
	defer resp.Body.Close()

	logger.Printf("%s GET: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	// on initial run there won't be any state
	if resp.StatusCode == http.StatusNoContent {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response code from GET to %s: %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "error reading response from GET to %s", url)
	}

	return string(body), nil

}

func saveLastID(key string, id string) error {

	d := &StateData{
		Key:   key,
		Value: id,
		Options: &StateOptions{
			Consistency: "strong",
		},
	}
	s := []*StateData{d}
	b, _ := json.Marshal(s)
	req, err := http.NewRequest(http.MethodPost, stateURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error posting to %s with key: %s, id: %s", stateURL, key, id)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", stateURL, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusCreated {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with key: %s, id: %s - %q",
			stateURL, key, id, dump)
	}

	return nil

}
