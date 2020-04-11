package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingQueryKey(t *testing.T) {

	q := &Query{
		Text:     "dapr",
		Lang:     "en",
		Count:    100,
		Username: "test",
		Token:    "test",
		Secret:   "test",
	}
	key1, err := parseQueryKey(q)
	assert.Nil(t, err)
	assert.NotEmpty(t, key1)

	key2, err := parseQueryKey(q)
	assert.Nil(t, err)
	assert.NotEmpty(t, key2)

	assert.Equal(t, key1, key2)

}

// func TestQueryHandler(t *testing.T) {

// 	gin.SetMode(gin.ReleaseMode)

// 	r := gin.Default()
// 	r.POST("/query", queryHandler)
// 	w := httptest.NewRecorder()

// 	q := &Query{
// 		Text:     "dapr",
// 		Lang:     "en",
// 		Count:    100,
// 		Username: "test",
// 		Token:    "test",
// 		Secret:   "test",
// 	}
// 	data, err := json.Marshal(q)
// 	assert.Nil(t, err)

// 	req, _ := http.NewRequest("POST", "/query", bytes.NewBuffer(data))
// 	req.Header.Set("Content-Type", "application/json")

// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// }
