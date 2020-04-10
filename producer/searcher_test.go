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

	c := &Config{
		Key:    "test",
		Secret: "test",
	}

	list, err := search(c, q)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, len(list), count)

}

func TestSearchErrors(t *testing.T) {

	_, err := search(nil, nil)
	assert.NotNil(t, err)

	c := &Config{
		Key: "test",
	}

	_, err = search(c, nil)
	assert.NotNil(t, err)

	tooHighCount := maxTweets + 1
	q := &Query{
		Text:     "test",
		Username: "test",
		Token:    "test",
		Secret:   "test",
		Count:    tooHighCount,
	}

	_, err = search(c, q)
	assert.NotNil(t, err)

}
