package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSearch(t *testing.T) {

	daprClient = GetTestClient()

	count := 50

	q := &Query{
		Query: "test",
		Count: count,
	}

	list, err := search(q)
	assert.Nil(t, err)
	assert.NotNil(t, list)

}

func TestSearchErrors(t *testing.T) {

	daprClient = GetTestClient()

	_, err := search(nil)
	assert.NotNil(t, err)

}
