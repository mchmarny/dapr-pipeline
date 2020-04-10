package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      serviceVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func mockHandler(c *gin.Context) {

	query := &Query{
		Text:     "dapr",
		Username: "dapr",
		Token:    "mock",
		Secret:   "mock",
	}

	c.JSON(http.StatusOK, query)

}

func queryHandler(c *gin.Context) {

	var q Query
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "unable to parse query",
		})
		return
	}

	list, err := search(queryConfig, &q)
	if err != nil {
		logger.Printf("error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error error executing query",
			"status":  "Failure",
		})
		return
	}

	c.JSON(http.StatusOK, list)
}
