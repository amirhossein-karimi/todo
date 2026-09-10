package middleware

import (
	"strconv"
	"time"

	"github.com/amirhossein-karimi/todo/internal/metrics"
	"github.com/gin-gonic/gin"
)

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		path := c.FullPath()

		if path == "" {
			path = "unknown"
		}

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		metrics.RequestsTotal.WithLabelValues(
			method,
			path,
			status,
		).Inc()

		metrics.RequestLatency.WithLabelValues(
			method,
			path,
		).Observe(duration)
	}
}
