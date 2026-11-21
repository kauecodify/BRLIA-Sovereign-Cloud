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
	"github.com/brlia/control-api/internal/config"
	"github.com/brlia/control-api/internal/handlers"
	"github.com/brlia/control-api/internal/middleware"
	"github.com/brlia/control-api/internal/repository"
	"github.com/brlia/control-api/internal/service"
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

	// Initialize repositories
	orgRepo := repository.NewOrganizationRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	quotaRepo := repository.NewQuotaRepository(db)
	regionRepo := repository.NewRegionRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// Initialize services
	orgService := service.NewOrganizationService(orgRepo, auditRepo)
	projectService := service.NewProjectService(projectRepo, orgRepo, auditRepo)
	userService := service.NewUserService(userRepo, auditRepo)
	roleService := service.NewRoleService(roleRepo, auditRepo)
	quotaService := service.NewQuotaService(quotaRepo, orgRepo, projectRepo, auditRepo)
	regionService := service.NewRegionService(regionRepo, auditRepo)
	auditService := service.NewAuditService(auditRepo)

	// Initialize HTTP server
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Format: "[BRLIA-CONTROL] '${time_stamp}' '${method}' '${path}' '${status}' '${latency}' '${error}'\n",
		TimestampFormat: time.RFC3339,
	}))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// Health check
	router.GET("/health", handlers.HealthCheckHandler())
	router.GET("/ready", handlers.ReadyCheckHandler(db, redisClient))

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Public routes
		apiV1.POST("/auth/token", handlers.GetTokenHandler(cfg.Auth))
		apiV1.POST("/auth/refresh", handlers.RefreshTokenHandler(cfg.Auth))

		// Authenticated routes
		authGroup := apiV1.Group("/")
		if cfg.Auth.Enabled {
			authGroup.Use(middleware.AuthMiddleware(cfg.Auth))
		}
		{
			// Organizations
			authGroup.GET("/organizations", handlers.ListOrganizationsHandler(orgService))
			authGroup.POST("/organizations", handlers.CreateOrganizationHandler(orgService))
			authGroup.GET("/organizations/:id", handlers.GetOrganizationHandler(orgService))
			authGroup.PUT("/organizations/:id", handlers.UpdateOrganizationHandler(orgService))
			authGroup.DELETE("/organizations/:id", handlers.DeleteOrganizationHandler(orgService))

			// Projects
			authGroup.GET("/projects", handlers.ListProjectsHandler(projectService))
			authGroup.POST("/projects", handlers.CreateProjectHandler(projectService))
			authGroup.GET("/projects/:id", handlers.GetProjectHandler(projectService))
			authGroup.PUT("/projects/:id", handlers.UpdateProjectHandler(projectService))
			authGroup.DELETE("/projects/:id", handlers.DeleteProjectHandler(projectService))

			// Users
			authGroup.GET("/users", handlers.ListUsersHandler(userService))
			authGroup.POST("/users", handlers.CreateUserHandler(userService))
			authGroup.GET("/users/:id", handlers.GetUserHandler(userService))
			authGroup.PUT("/users/:id", handlers.UpdateUserHandler(userService))
			authGroup.DELETE("/users/:id", handlers.DeleteUserHandler(userService))

			// Roles
			authGroup.GET("/roles", handlers.ListRolesHandler(roleService))
			authGroup.POST("/roles", handlers.CreateRoleHandler(roleService))
			authGroup.GET("/roles/:id", handlers.GetRoleHandler(roleService))
			authGroup.PUT("/roles/:id", handlers.UpdateRoleHandler(roleService))
			authGroup.DELETE("/roles/:id", handlers.DeleteRoleHandler(roleService))

			// Quotas
			authGroup.GET("/quotas", handlers.ListQuotasHandler(quotaService))
			authGroup.POST("/quotas", handlers.CreateQuotaHandler(quotaService))
			authGroup.GET("/quotas/:id", handlers.GetQuotaHandler(quotaService))
			authGroup.PUT("/quotas/:id", handlers.UpdateQuotaHandler(quotaService))

			// Regions
			authGroup.GET("/regions", handlers.ListRegionsHandler(regionService))
			authGroup.POST("/regions", handlers.CreateRegionHandler(regionService))
			authGroup.GET("/regions/:id", handlers.GetRegionHandler(regionService))

			// Audit
			authGroup.GET("/audit/logs", handlers.ListAuditLogsHandler(auditService))
			authGroup.GET("/audit/logs/:id", handlers.GetAuditLogHandler(auditService))

			// Resources
			authGroup.GET("/resources", handlers.ListResourcesHandler())
			authGroup.GET("/resources/:type", handlers.GetResourcesByTypeHandler())

			// Costs
			authGroup.GET("/costs", handlers.GetCostsHandler())
			authGroup.GET("/costs/summary", handlers.GetCostSummaryHandler())
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

	// Shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
