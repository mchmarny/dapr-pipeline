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
	logger = log.New(os.Stdout, "PROVIDER == ", 0)

	// service
	servicePort    = env.MustGetEnvVar("PORT", "8080")
	serviceVersion = env.MustGetEnvVar("RELEASE", "v0.0.1-default")

	// twitter
	consumerKey    = env.MustGetEnvVar("TW_CONSUMER_KEY", "")
	consumerSecret = env.MustGetEnvVar("TW_CONSUMER_SECRET", "")

	// dapr
	daprServer = fmt.Sprintf("http://localhost:%s", env.MustGetEnvVar("DAPR_HTTP_PORT", "3500"))
	daprClient = dapr.NewClient(daprServer)

	stateStore = env.MustGetEnvVar("STATE_STORE_NAME", "statestore")
	eventTopic = env.MustGetEnvVar("EVENT_TOPIC_NAME", "messagebus")
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
