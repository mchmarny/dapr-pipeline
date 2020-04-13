module github.com/mchmarny/dapr-pipeline/processor

go 1.14

require (
	github.com/cdipaolo/goml v0.0.0-20190412180403-e1f51f713598 // indirect
	github.com/cdipaolo/sentiment v0.0.0-20170111084539-5ab0aec020b4
	github.com/cloudevents/sdk-go/v2 v2.0.0-preview8
	github.com/gin-gonic/gin v1.6.2
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.0 // indirect
	github.com/mchmarny/gcputil v0.3.3
	github.com/mchmarny/godapr v0.2.3
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.22.3 // indirect
	go.uber.org/zap v1.14.1 // indirect
	golang.org/x/sys v0.0.0-20200413165638-669c56c373c4 // indirect
)

replace github.com/mchmarny/godapr => ../../godapr
