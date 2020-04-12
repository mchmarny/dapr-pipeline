package main

import (
	"github.com/cdipaolo/sentiment"
)

const (
	positiveSentimentScore = 1
	negativeSentimentScore = 0
)

var (
	model *sentiment.Models
)

func initModel() {
	m, err := sentiment.Restore()
	if err != nil {
		logger.Fatalf("error resting model: %v", err)
	}
	model = &m
}

func score(txt string) int {

	if model == nil {
		initModel()
	}

	a := model.SentimentAnalysis(txt, sentiment.English)
	// logger.Printf("sentiment (txt:'%s' == %+v)", txt, a)

	return int(a.Score)

}
