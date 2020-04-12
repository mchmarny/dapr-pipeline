module github.com/mchmarny/dapr-tweet-processing-pipeline/processor

go 1.14

replace github.com/mchmarny/dapr-tweet-processing-pipeline/client => ../client

require (
	github.com/gin-gonic/gin v1.6.2
	github.com/mchmarny/dapr-tweet-processing-pipeline/client v0.0.0-00010101000000-000000000000
	github.com/mchmarny/gcputil v0.3.3
	github.com/stretchr/testify v1.4.0
)
