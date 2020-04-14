package main

import (
	"fmt"
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
	daprServer = fmt.Sprintf("http://localhost:%s", env.MustGetEnvVar("DAPR_HTTP_PORT", "3500"))
	daprClient = dapr.NewClient(daprServer)

	stateStore = env.MustGetEnvVar("PRODUCER_STATE_STORE_NAME", "producer")
	eventTopic = env.MustGetEnvVar("PRODUCER_RESULT_TOPIC_NAME", "tweets")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	// simple routes
	r.GET("/", defaultHandler)
	r.POST("/query", queryHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server starting: %s \n", hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}
