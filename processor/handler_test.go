package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/", subscribeHandler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var got []string
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)

}

func TestEventProcessorHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/", eventHandler)
	w := httptest.NewRecorder()

	msg := gin.H{
		"k1": "v1",
		"k2": "v2",
		"ts": time.Now(),
	}

	data, err := json.Marshal(msg)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
