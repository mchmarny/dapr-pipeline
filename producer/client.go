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

// StateData represents ID state item in dapr
type StateData struct {
	Key     string        `json:"key"`
	Value   interface{}   `json:"value"`
	Options *StateOptions `json:"options,omitempty"`
}

// StateOptions is the dapr state data option for StateData
type StateOptions struct {
	Consistency string `json:"consistency,omitempty"`
}

func getState(key string) (data []byte, err error) {

	url := fmt.Sprintf("%s/%s", stateURL, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("consistency", "strong")

	client := &http.Client{
		Timeout: clientTimeout,
	}

	resp, err := client.Do(req)
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

func saveState(key string, data interface{}) error {

	state := &StateData{
		Key:     key,
		Value:   data,
		Options: &StateOptions{Consistency: "strong"},
	}
	list := []*StateData{state}
	b, _ := json.Marshal(list)
	req, err := http.NewRequest(http.MethodPost, stateURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error posting to %s with key: %s, data: %v", stateURL, key, data)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", stateURL, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusCreated {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with key: %s, data: %v - %q",
			stateURL, key, data, dump)
	}

	return nil

}

func publish(data interface{}) error {

	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, busURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error publishing result %+v to %s", data, busURL)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", busURL, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with result: %+v - %q",
			busURL, data, dump)
	}

	return nil

}
