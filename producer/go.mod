module github.com/mchmarny/dapr-tweet-processing-pipeline/producer

go 1.14

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dghubble/go-twitter v0.0.0-20190719072343-39e5462e111f
	github.com/dghubble/oauth1 v0.6.0
	github.com/gin-gonic/gin v1.6.2
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/mchmarny/dapr-tweet-processing-pipeline/client v0.0.0-00010101000000-000000000000
	github.com/mchmarny/gcputil v0.3.3
	github.com/stretchr/testify v1.5.1
)

replace github.com/mchmarny/dapr-tweet-processing-pipeline/client => ../client
