package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
)

const (
	// SupportedCloudEventVersion indicates the version of CloudEvents suppored by this handler
	SupportedCloudEventVersion = "0.3"

	//SupportedCloudEventContentTye indicates the content type supported by this handlers
	SupportedCloudEventContentTye = "application/json"

	sentimentAlertThreshold = float64(0.3)
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      AppVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func subscribeHandler(c *gin.Context) {
	topics := []string{sourceTopic}
	c.JSON(http.StatusOK, topics)
}

func eventHandler(c *gin.Context) {

	e := ce.NewEvent()
	if err := c.ShouldBindJSON(&e); err != nil {
		logger.Printf("error binding event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Bad Request",
			"message": "Error processing your request, see logs for details",
		})
		return
	}

	// logger.Printf("received event: %v", e.Context)

	eventVersion := e.Context.GetSpecVersion()
	if eventVersion != SupportedCloudEventVersion {
		logger.Printf("invalid event spec version: %s", eventVersion)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bad Request",
			"message": fmt.Sprintf("Invalid spec version (want: %s got: %s)",
				SupportedCloudEventVersion, eventVersion),
		})
		return
	}

	eventContentType := e.Context.GetDataContentType()
	if eventContentType != SupportedCloudEventContentTye {
		logger.Printf("invalid event content type: %s", eventContentType)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bad Request",
			"message": fmt.Sprintf("Invalid content type (want: %s got: %s)",
				SupportedCloudEventContentTye, eventContentType),
		})
		return
	}

	var t SimpleContent
	if err := json.Unmarshal(e.Data(), &t); err != nil {
		logger.Printf("error parsing event content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Bad Request",
			"message": "Invalid content payload, see log processor log for details",
		})
		return
	}

	// logger.Printf("tweet: %v", t)

	// score the content sentiment
	sentiment, err := score(t.Content)
	if err != nil {
		logger.Printf("error scoring sentiment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Server Error",
			"message": "Error scoring sentiment, see processor log for details",
		})
		return
	}
	t.Sentiment = sentiment

	// publish all
	if err := daprClient.Publish(processedTopic, t); err != nil {
		logger.Printf("error on processor result publish %v: %v", t, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Server Error",
			"message": "Error publishing result, see processor log for details",
		})
		return
	}

	// if negative then send alert
	if t.Sentiment <= sentimentAlertThreshold {
		logger.Printf("alert threshold %f reached: %f, sending alert", sentimentAlertThreshold, t.Sentiment)
		if err := daprClient.Send(alertBinding, t); err != nil {
			logger.Printf("error on sendign alert to binding %v: %v", alertBinding, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Server Error",
				"message": "Error sending alert, see processor log for details",
			})
			return
		}
	}

	c.JSON(http.StatusOK, nil)
}

// SimpleContent represents most
type SimpleContent struct {
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
	// Sentiment is the 0 to 1 score of the sentiment
	// (0-0.3 bad, 0.3-0.6 neutral, 0.6-1 positive)
	Sentiment float64 `json:"sentiment"`
	// Published is the parsed tweet create timestamp
	Published time.Time `json:"published"`
}
