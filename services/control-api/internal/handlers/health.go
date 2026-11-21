package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthCheckHandler returns a simple health check response
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "brlia-control-api",
			"timestamp": c.GetTime("timestamp"),
		})
	}
}

// ReadyCheckHandler returns a readiness check response
func ReadyCheckHandler(db *pgxpool.Pool, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"service":   "brlia-control-api",
				"database":  "unavailable",
				"error":     err.Error(),
				"timestamp": c.GetTime("timestamp"),
			})
			return
		}

		// Check Redis connection
		if _, err := redisClient.Ping(c.Request.Context()).Result(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"service":   "brlia-control-api",
				"database":  "healthy",
				"redis":     "unavailable",
				"error":     err.Error(),
				"timestamp": c.GetTime("timestamp"),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"service":   "brlia-control-api",
			"database":  "healthy",
			"redis":     "healthy",
			"timestamp": c.GetTime("timestamp"),
		})
	}
}
