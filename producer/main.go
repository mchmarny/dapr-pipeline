package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// service
	servicePort    = env.MustGetEnvVar("PORT", "8080")
	serviceVersion = env.MustGetEnvVar("RELEASE", "v0.0.1-default")

	// twitter
	queryConfig = &Config{
		Key:    env.MustGetEnvVar("TW_CONSUMER_KEY", ""),
		Secret: env.MustGetEnvVar("TW_CONSUMER_SECRET", ""),
	}

	// dapr
	daprPort = env.MustGetEnvVar("DAPR_HTTP_PORT", "3500")
	stateURL = fmt.Sprintf("http://localhost:%s/v1.0/state/%s", daprPort,
		env.MustGetEnvVar("DAPR_STORE_NAME", "statestore"))
	queueURL = fmt.Sprintf("http://localhost:%s/v1.0/publish/%s", daprPort,
		env.MustGetEnvVar("DAPR_QUEUE_NAME", "messagebus"))
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
