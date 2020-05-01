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

func TestParsingQueryKey(t *testing.T) {

	daprClient = GetTestClient()

	q := &Query{
		Query: "dapr",
		Lang:  "en",
		Count: 100,
	}
	key1, err := parseQueryKey(q)
	assert.Nil(t, err)
	assert.NotEmpty(t, key1)

	key2, err := parseQueryKey(q)
	assert.Nil(t, err)
	assert.NotEmpty(t, key2)

	assert.Equal(t, key1, key2)

}

func TestQueryHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/query", queryHandler)
	w := httptest.NewRecorder()

	q := &Query{
		Query: "dapr",
		Lang:  "en",
		Count: 100,
	}
	data, err := json.Marshal(q)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/query", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

func GetTestClient() *TestClient {
	return &TestClient{
		data: make(map[string][]byte, 0),
	}
}

type TestClient struct {
	data map[string][]byte
}

func (c *TestClient) GetData(store, key string) (data []byte, err error) {
	return c.data[key], nil
}
func (c *TestClient) SaveData(store, key string, data interface{}) error {
	b, _ := json.Marshal(data)
	c.data[key] = b
	return nil
}
func (c *TestClient) Publish(topic string, data interface{}) error {
	return nil
}
