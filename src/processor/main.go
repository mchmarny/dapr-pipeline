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
	logger = log.New(os.Stdout, "PROCESSOR == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8082")

	// dapr
	daprClient Client

	sourceTopic    = env.MustGetEnvVar("PROCESSOR_SOURCE_TOPIC_NAME", "tweets")
	processedTopic = env.MustGetEnvVar("PROCESSOR_RESULT_TOPIC_NAME", "processed")
	alertBinding   = env.MustGetEnvVar("PROCESSOR_ALERT_BINDING_NAME", "alert")

	apiEndpoint = env.MustGetEnvVar("PROCESSOR_API_ENDPOINT", "westus2.api.cognitive.microsoft.com")
	apiToken    = env.MustGetEnvVar("PROCESSOR_API_TOKEN", "")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// client
	daprClient = dapr.NewClient()

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
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}

type Client interface {
	Publish(topic string, data interface{}) error
	Send(binding string, data interface{}) error
}
