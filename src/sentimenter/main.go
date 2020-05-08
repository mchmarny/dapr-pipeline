package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var (
	logger = log.New(os.Stdout, "SENTIMENTER == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8082")

	apiEndpoint = env.MustGetEnvVar("SENTIMENTER_API_ENDPOINT", "westus2.api.cognitive.microsoft.com")
	apiToken    = env.MustGetEnvVar("CS_TOKEN", "")
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

	// simple routes
	r.GET("/", defaultHandler)
	r.POST("/score", scoreHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := http.ListenAndServe(hostPort, &ochttp.Handler{Handler: r}); err != nil {
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
