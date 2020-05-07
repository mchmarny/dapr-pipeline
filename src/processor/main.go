package main

import (
	"log"
	"net"
	"net/http"
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

	stateStore   = env.MustGetEnvVar("PROCESSOR_STATE_STORE_NAME", "tweet-store")
	eventTopic   = env.MustGetEnvVar("PROCESSOR_PUBSUB_TOPIC_NAME", "processed")
	scoreService = env.MustGetEnvVar("PROCESSOR_SCORE_SERVICE_NAME", "sentimenter")
	scoreMethod  = env.MustGetEnvVar("PROCESSOR_SCORE_METHOD_NAME", "score")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	daprClient = dapr.NewClient()

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

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

// Options midleware
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Client is the minim client support for testing
type Client interface {
	SaveState(store, key string, data interface{}) error
	InvokeService(service, method string, data interface{}) (out []byte, err error)
	Publish(topic string, data interface{}) error
}
