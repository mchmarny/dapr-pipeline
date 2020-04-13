module github.com/mchmarny/dapr-pipeline/producer

go 1.14

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dghubble/go-twitter v0.0.0-20190719072343-39e5462e111f
	github.com/dghubble/oauth1 v0.6.0
	github.com/gin-gonic/gin v1.6.2
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/mchmarny/gcputil v0.3.3
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/sys v0.0.0-20200413165638-669c56c373c4 // indirect
)

replace github.com/mchmarny/godapr => ../../godapr
