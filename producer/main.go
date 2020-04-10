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

	servicePort    = env.MustGetEnvVar("PORT", "8080")
	serviceVersion = env.MustGetEnvVar("RELEASE", "v0.0.1-default")

	consumerKey    = env.MustGetEnvVar("TW_CONSUMER_KEY", "")
	consumerSecret = env.MustGetEnvVar("TW_CONSUMER_SECRET", "")

	daprPort  = env.MustGetEnvVar("DAPR_HTTP_PORT", "3500")
	daprStore = env.MustGetEnvVar("DAPR_STORE", "tweets")
	storeURL  = fmt.Sprintf("http://localhost:%s/v1.0/state/%s", daprPort, daprStore)
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	// simple routes
	r.GET("/", defaultHandler)

	// api
	v1 := r.Group("/v1")
	{
		v1.POST("/notif", queryHandler)
	}

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server starting: %s \n", hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}

}
