package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSearch(t *testing.T) {

	q := &Query{
		Text:     "test",
		Username: "test",
		Token:    "test",
		Secret:   "test",
	}

	c := &Config{
		Key:    "test",
		Secret: "test",
	}

	list, err := search(c, q)
	assert.Nil(t, err)
	assert.NotNil(t, list)

}
