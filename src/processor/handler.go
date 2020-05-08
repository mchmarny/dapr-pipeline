package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dghubble/go-twitter/twitter"
)

var (
	clientError = gin.H{
		"error":   "Bad Request",
		"message": "Error processing your request, see logs for details",
	}
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      AppVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func tweetHandler(c *gin.Context) {
	var t twitter.Tweet
	if err := c.ShouldBindJSON(&t); err != nil {
		logger.Printf("error binding tweet: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("tweet: %s", t.IDStr)

	ctx := c.Request.Context()

	// save original tweet in case we need to reprocess it
	err := daprClient.SaveState(ctx, stateStore, t.IDStr, t)
	if err != nil {
		logger.Printf("error saving state: %v", err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	content := t.FullText
	if content == "" {
		content = t.Text
	}

	sentimentReq := struct {
		Text string `json:"text"`
		Lang string `json:"lang"`
	}{
		content,
		t.Lang,
	}

	// score simple tweet
	b, err := daprClient.InvokeService(ctx, scoreService, scoreMethod, sentimentReq)
	if err != nil {
		logger.Printf("error invoking scoring service (%s/%s): %v",
			scoreService, scoreMethod, err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	sentimentRes := struct {
		Score float64 `json:"score"`
	}{}

	if err := json.Unmarshal(b, &sentimentRes); err != nil {
		logger.Printf("error parsing scoring service response (%s): %v", string(b), err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	// create simple tweet from status for scoring and display
	s := &SimpleTweet{
		ID:        t.IDStr,
		Author:    strings.ToLower(t.User.ScreenName),
		AuthorPic: t.User.ProfileImageURLHttps,
		Published: convertTwitterTime(t.CreatedAt),
		Content:   content,
		Score:     sentimentRes.Score,
	}

	// publish simple tweet
	if err = daprClient.Publish(ctx, eventTopic, s); err != nil {
		logger.Printf("error publishing content (%+v): %v", s, err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func convertTwitterTime(v string) time.Time {
	t, err := time.Parse(time.RubyDate, v)
	if err != nil {
		t = time.Now()
	}
	return t.UTC()
}

// SimpleTweet represents the Twiter query result item
type SimpleTweet struct {
	// ID is the string representation of the tweet ID
	ID string `json:"id"`
	// Author is the name of the tweet user
	Author string `json:"author"`
	// AuthorPic is the url to author profile pic
	AuthorPic string `json:"author_pic"`
	// Content is the full text body of the tweet
	Content string `json:"content"`
	// Content is the sentiment score of the content
	Score float64 `json:"score"`
	// Published is the parsed tweet create timestamp
	Published time.Time `json:"published"`
}
