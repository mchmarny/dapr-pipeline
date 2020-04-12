package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"

	dapr "github.com/mchmarny/dapr-tweet-processing-pipeline/client"
)

var (
	logger = log.New(os.Stdout, "PROCESSOR == ", 0)

	// service
	servicePort    = env.MustGetEnvVar("PORT", "8080")
	serviceVersion = env.MustGetEnvVar("RELEASE", "v0.0.1-default")

	// dapr
	daprServer = fmt.Sprintf("http://localhost:%s", env.MustGetEnvVar("DAPR_HTTP_PORT", "3500"))
	daprClient = dapr.NewClient(daprServer)

	sourceTopic   = env.MustGetEnvVar("PROCESSOR_SOURCE_TOPIC_NAME", "tweets")
	positiveTopic = env.MustGetEnvVar("PROCESSOR_POSITIVE_TOPIC_NAME", "positive-tweets")
	negativeTopic = env.MustGetEnvVar("PROCESSOR_NEGATIVE_TOPIC_NAME", "negative-tweets")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	// simple routes
	r.GET("/", defaultHandler)
	r.GET("/dapr/subscribe", subscribeHandler)

	// topic route
	processRoute := fmt.Sprintf("/%s", sourceTopic)
	logger.Printf("processor route: %s", processRoute)
	r.POST(processRoute, eventHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server starting: %s \n", hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}
