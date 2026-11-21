package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/brlia/waterfall/internal/config"
	"github.com/brlia/waterfall/internal/engine"
	"github.com/brlia/waterfall/internal/handlers"
	"github.com/brlia/waterfall/internal/middleware"
	"github.com/brlia/waterfall/internal/repository"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	// Initialize Redis connection
	redisClient, err := repository.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize Waterfall Engine
	waterfallEngine, err := engine.NewWaterfallEngine(cfg, db, redisClient)
	if err != nil {
		log.Fatalf("Failed to initialize Waterfall Engine: %v", err)
	}

	// Start background workers
	go waterfallEngine.StartWorkers()

	// Initialize HTTP server
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Format: "[BRLIA-WATERFALL] '${time_stamp}' '${method}' '${path}' '${status}' '${latency}' '${error}'\n",
		TimestampFormat: time.RFC3339,
	}))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// Health check
	router.GET("/health", handlers.HealthCheckHandler())
	router.GET("/ready", handlers.ReadyCheckHandler(db, redisClient, waterfallEngine))

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Public routes
		apiV1.GET("/status", handlers.GetEngineStatusHandler(waterfallEngine))

		// Authenticated routes
		authGroup := apiV1.Group("/")
		if cfg.Auth.Enabled {
			authGroup.Use(middleware.AuthMiddleware(cfg.Auth))
		}
		{
			// Workload management
			authGroup.POST("/workloads", handlers.CreateWorkloadHandler(waterfallEngine))
			authGroup.GET("/workloads", handlers.ListWorkloadsHandler(waterfallEngine))
			authGroup.GET("/workloads/:id", handlers.GetWorkloadHandler(waterfallEngine))
			authGroup.PUT("/workloads/:id", handlers.UpdateWorkloadHandler(waterfallEngine))
			authGroup.DELETE("/workloads/:id", handlers.DeleteWorkloadHandler(waterfallEngine))

			// Resource management
			authGroup.GET("/resources", handlers.ListResourcesHandler(waterfallEngine))
			authGroup.GET("/resources/:id", handlers.GetResourceHandler(waterfallEngine))

			// Capacity management
			authGroup.GET("/capacity", handlers.GetCapacityHandler(waterfallEngine))
			authGroup.GET("/capacity/forecast", handlers.GetCapacityForecastHandler(waterfallEngine))

			// Queue management
			authGroup.GET("/queues", handlers.ListQueuesHandler(waterfallEngine))
			authGroup.GET("/queues/:name", handlers.GetQueueHandler(waterfallEngine))

			// Scheduling
			authGroup.POST("/schedule", handlers.ScheduleWorkloadHandler(waterfallEngine))
			authGroup.GET("/schedule", handlers.GetScheduleHandler(waterfallEngine))

			// Policies
			authGroup.GET("/policies", handlers.ListPoliciesHandler(waterfallEngine))
			authGroup.POST("/policies", handlers.CreatePolicyHandler(waterfallEngine))
			authGroup.GET("/policies/:id", handlers.GetPolicyHandler(waterfallEngine))
			authGroup.PUT("/policies/:id", handlers.UpdatePolicyHandler(waterfallEngine))
			authGroup.DELETE("/policies/:id", handlers.DeletePolicyHandler(waterfallEngine))

			// Metrics
			authGroup.GET("/metrics", handlers.GetMetricsHandler(waterfallEngine))
			authGroup.GET("/metrics/workloads", handlers.GetWorkloadMetricsHandler(waterfallEngine))
			authGroup.GET("/metrics/resources", handlers.GetResourceMetricsHandler(waterfallEngine))

			// Cost optimization
			authGroup.GET("/costs/optimization", handlers.GetCostOptimizationHandler(waterfallEngine))
			authGroup.POST("/costs/analyze", handlers.AnalyzeCostsHandler(waterfallEngine))

			// Failure detection and recovery
			authGroup.GET("/failures", handlers.ListFailuresHandler(waterfallEngine))
			authGroup.GET("/failures/:id", handlers.GetFailureHandler(waterfallEngine))
			authGroup.POST("/recover", handlers.TriggerRecoveryHandler(waterfallEngine))

			// Multi-cluster orchestration
			authGroup.GET("/clusters", handlers.ListClustersHandler(waterfallEngine))
			authGroup.GET("/clusters/:id", handlers.GetClusterHandler(waterfallEngine))
			authGroup.POST("/clusters/:id/workloads", handlers.AssignWorkloadToClusterHandler(waterfallEngine))
		}
	}

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop workers
	waterfallEngine.StopWorkers()

	// Shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
