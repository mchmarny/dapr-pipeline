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
	logger = log.New(os.Stdout, "PROVIDER == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8081")

	// twitter
	consumerKey    = env.MustGetEnvVar("TW_CONSUMER_KEY", "")
	consumerSecret = env.MustGetEnvVar("TW_CONSUMER_SECRET", "")
	accessToken    = env.MustGetEnvVar("TW_ACCESS_TOKEN", "")
	accessSecret   = env.MustGetEnvVar("TW_ACCESS_SECRET", "")

	// dapr
	daprClient Client

	stateStore = env.MustGetEnvVar("PRODUCER_STATE_STORE_NAME", "producer")
	eventTopic = env.MustGetEnvVar("PRODUCER_RESULT_TOPIC_NAME", "tweets")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	daprClient = dapr.NewClient()

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	// simple routes
	r.GET("/", defaultHandler)
	r.POST("/query", queryHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}

type Client interface {
	GetData(store, key string) (data []byte, err error)
	SaveData(store, key string, data interface{}) error
	Publish(topic string, data interface{}) error
}
