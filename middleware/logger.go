package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a structured request logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		logLine := fmt.Sprintf("[%s] %s | %d | %v | %s | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			statusCode,
			latency,
			clientIP,
			path,
		)

		// Color-code by status
		switch {
		case statusCode >= 500:
			fmt.Printf("\033[31m%s\033[0m\n", logLine) // Red
		case statusCode >= 400:
			fmt.Printf("\033[33m%s\033[0m\n", logLine) // Yellow
		default:
			fmt.Printf("\033[32m%s\033[0m\n", logLine) // Green
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				fmt.Printf("\033[31m[ERROR] %s\033[0m\n", e.Error())
			}
		}
	}
}
