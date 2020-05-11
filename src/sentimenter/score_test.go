package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

const (
	goodText     = "I like how this food tastes, it makes me happy"
	negativeText = "Your football team is really bad, they are awful this season"
	lang         = "en"
)

// go test -v -count=1 -run TestScoring ./...
func TestScoring(t *testing.T) {
	ctx := trace.SpanContext{}

	gs, err := scoreSentiment(ctx, goodText, lang)
	assert.Nil(t, err, "error scoring good text")
	assert.GreaterOrEqual(t, gs, float64(0.6))

	bs, err := scoreSentiment(ctx, negativeText, lang)
	assert.Nil(t, err, "error scoring negative text")
	assert.LessOrEqual(t, bs, float64(0.3))
}
