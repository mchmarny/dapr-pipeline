package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	clientTimeout = time.Second * 30
)

// IDStateData represents ID state item in dapr
type IDStateData struct {
	Key     string          `json:"key"`
	Value   int64           `json:"value"`
	Options *IDStateOptions `json:"options,omitempty"`
}

// IDStateOptions is the dapr state data option for IDStateData
type IDStateOptions struct {
	Consistency string `json:"consistency,omitempty"`
}

func getLastID(key string) (id int64, err error) {

	url := fmt.Sprintf("%s/%s", stateURL, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("consistency", "strong")

	client := &http.Client{
		Timeout: clientTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.Wrapf(err, "error quering state service: %s", url)
	}
	defer resp.Body.Close()

	logger.Printf("%s GET: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	// on initial run there won't be any state
	if resp.StatusCode == http.StatusNoContent {
		return 0, nil
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid response code from GET to %s: %d", url, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading response from GET to %s", url)
	}

	idStr := strings.ReplaceAll(string(data), "\"", "") //HUCK: save as json object so can parse here
	lastID, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "error parsing response '%s' from GET to %s", idStr, url)
	}

	return lastID, nil

}

func saveLastID(key string, id int64) error {

	d := &IDStateData{
		Key:   key,
		Value: id,
		Options: &IDStateOptions{
			Consistency: "strong",
		},
	}
	s := []*IDStateData{d}
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
		return fmt.Errorf("invalid response code from POST to %s with key: %s, id: %d - %q",
			stateURL, key, id, dump)
	}

	return nil

}

func publishQueryResult(t *SimpleTweet) error {

	b, _ := json.Marshal(t)
	req, err := http.NewRequest(http.MethodPost, busURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error publishing result %+v to %s", t, busURL)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", busURL, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		dump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("invalid response code from POST to %s with result: %+v - %q",
			busURL, t, dump)
	}

	return nil

}
