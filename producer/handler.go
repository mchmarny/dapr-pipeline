package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
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

func queryHandler(c *gin.Context) {

	var q Query
	if err := c.ShouldBindJSON(&q); err != nil {
		logger.Printf("error binding query: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	// key
	queryKey, err := parseQueryKey(&q)
	if err != nil {
		logger.Printf("unable parse query key: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	if q.SinceID == "" {
		lastID, err := getLastID(queryKey)
		if err != nil {
			logger.Printf("error retrieving last query ID: %v", err)
			c.JSON(http.StatusBadRequest, clientError)
			return
		}
		if lastID != "" {
			// update the since ID with the one retrieved from dapr state
			logger.Printf("found last ID: %s", lastID)
			q.SinceID = lastID
		}
	}

	list, err := search(queryConfig, &q)
	if err != nil {
		logger.Printf("error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	// save last tweet ID in state if there are results
	if len(list) > 0 {
		newLastID := list[len(list)-1].ID
		logger.Printf("new last ID: %s", newLastID)
		err = saveLastID(queryKey, newLastID)
		if err != nil {
			logger.Printf("error saving new last ID: %v", err)
			c.JSON(http.StatusInternalServerError, clientError)
			return
		}
	}

	c.JSON(http.StatusOK, list)
}

func parseQueryKey(q *Query) (key string, err error) {

	if q == nil {
		return "", errors.New("nil query pointer")
	}

	if q.Text == "" {
		return "", errors.New("empty query text")
	}

	if q.Username == "" {
		return "", errors.New("empty query username")
	}

	rawKey := fmt.Sprintf("u:%s|t:%s", q.Username, q.Text)
	logger.Printf("raw key: %s", rawKey)

	hashedKey := fmt.Sprintf("qk-%s", toMD5Hash(rawKey))
	logger.Printf("hashed key: %s", hashedKey)

	return hashedKey, nil

}

func toMD5Hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
