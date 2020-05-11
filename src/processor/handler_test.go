package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestTweetHandler(t *testing.T) {

	daprClient = GetTestClient()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/tweets", tweetHandler)
	w := httptest.NewRecorder()

	data, err := ioutil.ReadFile("./tweet.json")
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/tweets", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func GetTestClient() *TestClient {
	return &TestClient{}
}

var (
	// test test client against local interace
	_ = Client(&TestClient{})
)

type TestClient struct {
}

func (c *TestClient) SaveState(ctx trace.SpanContext, store, key string, data interface{}) error {
	return nil
}

func (c *TestClient) InvokeService(ctx trace.SpanContext, service, method string, data interface{}) (out []byte, err error) {
	return []byte("{\"score\": 0.1234556789}"), nil
}

func (c *TestClient) Publish(ctx trace.SpanContext, topic string, data interface{}) error {
	return nil
}
