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

const (
	defaultDocID = "1"
	apiURL       = "https://%s/text/analytics/v2.1/sentiment"
)

type RequestItem struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Text     string `json:"text"`
}

type Request struct {
	Docs []*RequestItem `json:"documents"`
}

type ResponseItem struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
}

type Response struct {
	Docs []*ResponseItem `json:"documents"`
}

func score(txt string) (sentiment float64, err error) {
	if txt == "" {
		return 0, errors.New("nil txt")
	}

	// content array for request
	data := Request{
		Docs: []*RequestItem{
			&RequestItem{
				ID:       defaultDocID,
				Language: "en",
				Text:     txt,
			},
		},
	}

	b, _ := json.Marshal(data)
	url := fmt.Sprintf(apiURL, apiEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", apiToken)

	c := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, errors.Wrapf(err, "error posting %+v to %s", data, url)
	}
	defer resp.Body.Close()

	logger.Printf("%s POST: %d (%s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		dump, _ := httputil.DumpResponse(resp, true)
		return 0, errors.Wrapf(err, "invalid response code from POST to %s with result: %+v - %q",
			url, data, dump)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading response from POST to %s", url)
	}

	logger.Printf("DEBUG: %s", string(content))

	var r Response
	err = json.Unmarshal(content, &r)
	if err != nil {
		return 0, errors.Wrapf(err, "error parsing response from content: %s", string(content))
	}

	if r.Docs != nil && len(r.Docs) != 1 {
		return 0, errors.Wrapf(err, "expected 1 doc but API returned %d", len(r.Docs))
	}

	doc := r.Docs[0]
	if doc.ID != defaultDocID {
		return 0, errors.Wrapf(err, "expected ID %s but API returned %s", defaultDocID, doc.ID)
	}

	return doc.Score, nil
}
