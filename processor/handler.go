package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	clientError = gin.H{
		"error":   "Bad Request",
		"message": "Error processing your request, see logs for details",
	}
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      serviceVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func subscribeHandler(c *gin.Context) {
	topics := []string{sourceTopic}
	c.JSON(http.StatusOK, topics)
}

func eventHandler(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Printf("error reading request body: %v", err)
	}

	logger.Printf("%s", b)

	c.JSON(http.StatusOK, nil)
}
