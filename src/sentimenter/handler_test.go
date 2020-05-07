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

func TestScoreHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/score", scoreHandler)
	w := httptest.NewRecorder()

	in := &ScoreRequest{
		Text: "I'm so happy this test works",
	}

	data, _ := json.Marshal(in)

	req, _ := http.NewRequest("POST", "/score", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var out SimpleScore
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, in.Text, out.Text)

}
