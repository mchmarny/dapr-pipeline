package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

const (
	traceExporterNotSet = "trace-not-set"
)

var (
	logger = log.New(os.Stdout, "SENTIMENTER == ", 0)

	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	// service
	servicePort = env.MustGetEnvVar("PORT", "8082")

	apiEndpoint = env.MustGetEnvVar("SENTIMENTER_API_ENDPOINT", "westus2.api.cognitive.microsoft.com")
	apiToken    = env.MustGetEnvVar("CS_TOKEN", "")
	exporterURL = env.MustGetEnvVar("TRACE_EXPORTER_URL", traceExporterNotSet)
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// START TRACING
	if exporterURL != traceExporterNotSet {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "localhost"
		}
		endpointID := fmt.Sprintf("%s:%s", hostname, servicePort)
		localEndpoint, err := openzipkin.NewEndpoint("sentimenter", endpointID)
		if err != nil {
			logger.Fatalf("error creating local endpoint: %v", err)
		}
		reporter := zipkinHTTP.NewReporter(exporterURL)
		trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}
	// END TRACING

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
