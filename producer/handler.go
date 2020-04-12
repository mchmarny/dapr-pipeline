package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

	stateContent, err := daprClient.GetState(stateStore, queryKey)
	if err != nil {
		logger.Printf("error retrieving state: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	var lastMaxID int64
	if len(stateContent) > 0 {
		idStr := string(stateContent) //HUCK: save as json object so can parse here
		lastMaxID, err = strconv.ParseInt(idStr, 0, 64)
		if err != nil {
			logger.Printf("error parsing response '%s', using default %d: %v", idStr, lastMaxID, err)
		}
	}

	logger.Printf("found last max ID: %d", lastMaxID)
	q.SinceID = lastMaxID

	r, err := search(&q)
	if err != nil {
		logger.Printf("error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, clientError)
		return
	}

	logger.Printf("new maxID: %t (%d)", r.MaxID > lastMaxID, r.MaxID-lastMaxID)
	// save only if the new maxID is higher than the one already saved
	if r.MaxID > lastMaxID {
		err = daprClient.SaveState(stateStore, queryKey, r.MaxID)
		if err != nil {
			logger.Printf("error saving state: %v", err)
			c.JSON(http.StatusInternalServerError, clientError)
			return
		}
	}

	c.JSON(http.StatusOK, r)
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
	// logger.Printf("raw key: %s", rawKey)

	hashedKey := fmt.Sprintf("qk-%s", toMD5Hash(rawKey))
	// logger.Printf("hashed key: %s", hashedKey)

	return hashedKey, nil

}

func toMD5Hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
