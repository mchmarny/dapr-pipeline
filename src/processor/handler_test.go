package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	daprClient = &TestClient{}

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

	daprClient = &TestClient{}

	r := gin.Default()
	r.POST("/", eventHandler)
	w := httptest.NewRecorder()

	data := []byte(`{
		"id":"eeda6273-7483-4d4a-b368-41edbde76257",
		"source":"provider",
		"type":"com.dapr.event.sent",
		"specversion":"0.3",
		"datacontenttype":"application/json",
		"data":{
			"id":"1249435923240104999",
			"query":"serverless",
			"author":"test",
			"content":"test message",
			"published":"2020-04-12T20:35:12Z"
		},
		"subject":""
	}`)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

type TestClient struct {
}

func (c *TestClient) Publish(topic string, data interface{}) error {
	return nil
}

func (c *TestClient) Send(binding string, data interface{}) error {
	return nil
}
