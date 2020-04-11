package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	maxTweets = 100
)

// Query represents Twitter search query in specific user context
type Query struct {

	// Text is the full text of the Twitter search query including operators
	// e.g. 'dapr AND microsoft'
	Text string `json:"text"`

	// Lang is the ISO 639-1 code which will be used to filter tweets
	Lang string `json:"lang"`

	// Count is the number of tweets to return (no paging for now)
	Count int `json:"count"`

	// SinceID is the id of the tweet to start search from
	// If not provided it will be set to the last tweet returned by this query
	SinceID string `json:"since_id"`

	// Username is the Twitter username who's Token/Secrets are assciated with
	Username string `json:"user"`

	// Token is the Twitter AccessTokenKey
	Token string `json:"token"`

	// Secret is the Twitter AccessTokenSecrets
	Secret string `json:"secret"`
}

func (q *Query) validate() error {

	if q.Text == "" {
		return errors.New("empty search query text")
	}

	if q.Username == "" {
		return errors.New("empty search query user")
	}

	if q.Token == "" {
		return errors.New("empty search query token")
	}

	if q.Secret == "" {
		return errors.New("empty search query secret")
	}

	if q.Count == 0 {
		q.Count = maxTweets
	}

	if q.Count > 100 {
		logger.Printf("invalid query.count (want: 0-%d, got: %d), re-setting to max: %d",
			maxTweets, q.Count, maxTweets)
		q.Count = maxTweets
	}

	if q.Lang == "" {
		q.Lang = "en"
	}

	return nil

}

// Config is the application configuration information
type Config struct {
	// Token is the Twitter app config consumer key
	Key string `json:"token"`
	// Secret is the Twitter app config consumer secret
	Secret string `json:"secret"`
}

func (c *Config) validate() error {

	if c.Key == "" {
		return errors.New("empty search config key")
	}

	if c.Secret == "" {
		return errors.New("empty search config secret")
	}

	return nil

}

// Result represents the query result
type Result struct {
	// ID is the string representation of the tweet ID
	ID string `json:"id"`
	// Author is the name of the tweet user
	Author string `json:"author"`
	// Text is the full text body of the tweet
	Text string `json:"text"`
	// Published is the parsed tweet create timestamp
	Published time.Time `json:"published"`
}

func search(c *Config, q *Query) (r []*Result, err error) {

	if c == nil {
		return nil, errors.New("nil app config")
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	if q == nil {
		return nil, errors.New("nil search query")
	}

	if err := q.validate(); err != nil {
		return nil, err
	}

	r = make([]*Result, 0)

	// Mocking for now
	for i := 0; i < q.Count; i++ {
		item := &Result{
			ID:        strconv.Itoa(i),
			Author:    fmt.Sprintf("author%d", i),
			Text:      fmt.Sprintf("text %d", i),
			Published: time.Now().UTC(),
		}
		r = append(r, item)
	}

	return

}
