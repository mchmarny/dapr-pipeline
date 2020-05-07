package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"

	dapr "github.com/mchmarny/godapr"
)

var (
	logger = log.New(os.Stdout, "PROCESSOR == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8081")

	// dapr
	daprClient Client

	// test client against local interace
	_ = Client(dapr.NewClient())

	stateStore   = env.MustGetEnvVar("PRODUCER_STATE_STORE_NAME", "tweet-store")
	eventTopic   = env.MustGetEnvVar("PRODUCER_PUBSUB_TOPIC_NAME", "processed")
	scoreService = env.MustGetEnvVar("PRODUCER_SCORE_SERVICE_NAME", "sentimenter")
	scoreMethod  = env.MustGetEnvVar("PRODUCER_SCORE_METHOD_NAME", "score")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	daprClient = dapr.NewClient()

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	// simple routes
	r.GET("/", defaultHandler)
	r.POST("/tweets", tweetHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}

// Client is the minim client support for testing
type Client interface {
	SaveState(store, key string, data interface{}) error
	InvokeService(service, method string, data interface{}) (out []byte, err error)
	Publish(topic string, data interface{}) error
}
