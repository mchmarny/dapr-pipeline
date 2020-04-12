package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	goodText = "I like how this food tastes, it makes me happy"
	badText  = "Your team sucks and awful this season"
)

func TestScoring(t *testing.T) {
	assert.Equal(t, score(goodText), 1)
	assert.Equal(t, score(badText), 0)
}
