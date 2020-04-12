package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientCreation(t *testing.T) {
	url := "test"
	c := NewClient(url)
	assert.NotNil(t, c)
	assert.Equal(t, c.BaseURL, url)
}
