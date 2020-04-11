package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSearch(t *testing.T) {

	count := 50

	q := &Query{
		Text:     "test",
		Username: "test",
		Token:    "test",
		Secret:   "test",
		Count:    count,
	}

	list, err := search(q)
	assert.Nil(t, err)
	assert.NotNil(t, list)

}

func TestSearchErrors(t *testing.T) {

	_, err := search(nil)
	assert.NotNil(t, err)

	tooHighCount := maxTweets + 1
	q := &Query{
		Text:     "test",
		Username: "test",
		Token:    "test",
		Secret:   "test",
		Count:    tooHighCount,
	}

	_, err = search(q)
	assert.NotNil(t, err)

}
