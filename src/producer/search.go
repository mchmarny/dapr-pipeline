package main

import (
	"errors"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	maxTweets = 100
)

// Query represents Twitter search query in specific user context
type Query struct {

	// Query is the full text of the Twitter search query including operators
	// e.g. 'dapr AND microsoft'
	Query string `json:"query"`

	// Lang is the ISO 639-1 code which will be used to filter tweets
	Lang string `json:"lang"`

	// Count is the number of tweets to return (no paging for now)
	Count int `json:"count"`

	// SinceID is the id of the tweet to start search from
	// Set to the last tweet returned by this query in handler
	SinceID int64 `json:"-"`
}

func (q *Query) validate() error {

	if q.Query == "" {
		return errors.New("empty search query")
	}

	if q.Count == 0 {
		q.Count = maxTweets
	}

	if q.Count > 100 {
		logger.Printf("invalid query.count (want: 0-%d, got: %d), re-setting to max: %d",
			maxTweets, q.Count, maxTweets)
		q.Count = maxTweets
	}

	return nil

}

// SimpleTweet represents the Twiter query result item
type SimpleTweet struct {
	// ID is the string representation of the tweet ID
	ID string `json:"id"`
	// Query is the text of the original query
	Query string `json:"query"`
	// Author is the name of the tweet user
	Author string `json:"author"`
	// AuthorPic is the url to author profile pic
	AuthorPic string `json:"author_pic"`
	// Content is the full text body of the tweet
	Content string `json:"content"`
	// Published is the parsed tweet create timestamp
	Published time.Time `json:"published"`
}

// SearchResult is the metadata from executed search
type SearchResult struct {
	// LastID is the last tweet ID
	SinceID int64 `json:"since_id"`
	// MaxID is the last tweet ID
	MaxID int64 `json:"max_id"`
	// Query is the text of the search query
	Query string `json:"query"`
	// QueryKey is MD5 hash of the query
	QueryKey string `json:"query_key"`
	// Found is the number of items returned by search
	Found int `json:"items_found"`
	// Published is the number of items published
	Published int `json:"items_published"`
	// Duration is the number of items published
	Duration float64 `json:"search_duration"`
}

func search(q *Query) (r *SearchResult, err error) {

	if q == nil {
		return nil, errors.New("nil search query")
	}

	if err := q.validate(); err != nil {
		return nil, err
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	tc := twitter.NewClient(httpClient)

	// logger.Printf("searching for '%s' since id: %d", q.Text, q.SinceID)
	search, resp, err := tc.Search.Tweets(&twitter.SearchTweetParams{
		Query:      q.Query,
		Count:      maxTweets,
		Lang:       q.Lang,
		SinceID:    q.SinceID,
		ResultType: "recent",
		TweetMode:  "extended",
	})

	if err != nil {
		logger.Printf("error on search: %v - %v", resp, err)
		return nil, err
	}

	r = &SearchResult{
		Query:    search.Metadata.Query,
		SinceID:  q.SinceID,
		MaxID:    q.SinceID, // start with the previous max in case there is no more results
		Duration: search.Metadata.CompletedIn,
	}

	for _, s := range search.Statuses {

		// increment found count to comp to published later
		r.Found++

		// filter out RT
		// TODO: parameterize
		if s.RetweetedStatus != nil {
			// logger.Printf("skipping RT: %s", s.FullText)
			continue
		}

		// create simple tweet from status
		t := &SimpleTweet{
			ID:        s.IDStr,
			Query:     q.Query,
			Author:    strings.ToLower(s.User.ScreenName),
			AuthorPic: s.User.ProfileImageURLHttps,
			Published: convertTwitterTime(s.CreatedAt),
			Content:   s.FullText,
		}

		// publish simple tweet
		if err = daprClient.Publish(eventTopic, t); err != nil {
			logger.Printf("error on publish %v: %v", t, err)
			//return so we don't update the last ID and have chance to reporcess this query
			return nil, err
		}
		// set if current tweet ID larger than the current max
		// can't assume tweets arrive in latest last order
		if s.ID > r.MaxID {
			logger.Printf("new tweet, ID diff: %d", s.ID-r.MaxID)
			r.MaxID = s.ID
		}

		// increment published count
		r.Published++
	}

	return r, nil

}

func convertTwitterTime(v string) time.Time {
	t, err := time.Parse(time.RubyDate, v)
	if err != nil {
		t = time.Now()
	}
	return t.UTC()
}
