module github.com/mchmarny/dapr-tweet-processing-pipeline/processor

go 1.14

require (
	github.com/gin-gonic/gin v1.6.2
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/mchmarny/dapr-tweet-processing-pipeline/client v0.0.0-20200412211330-d1a4b0cca1fc
	github.com/mchmarny/gcputil v0.3.3
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/sys v0.0.0-20200409092240-59c9f1ba88fa // indirect
)

replace github.com/mchmarny/dapr-tweet-processing-pipeline/client => ../client
