package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      serviceVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func queryHandler(c *gin.Context) {

	// if e := publish(c.Request.Context(), data); e != nil {
	// 	logger.Printf("error publishing notification: %v", e)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Error handling notification",
	// 		"status":  "Failure",
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification proccessed",
		"status":  "OK",
	})
}
