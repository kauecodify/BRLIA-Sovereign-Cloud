package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/brlia/waterfall/internal/engine"
	"github.com/brlia/waterfall/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// store keeps an in-memory registry of workloads and policies
type store struct {
	mu        sync.RWMutex
	workloads map[uuid.UUID]*models.Workload
	policies  map[uuid.UUID]*models.Policy
}

var sharedStore = &store{
	workloads: make(map[uuid.UUID]*models.Workload),
	policies:  make(map[uuid.UUID]*models.Policy),
}

func init() {
	seedDemoPolicies()
}

func seedDemoPolicies() {
	now := time.Now().UTC()
	enabled := true
	demo := []models.Policy{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000101"),
			Name:        "prefer-national-regions",
			Description: "Prioriza alocacao em regioes nacionais soberanas",
			Type:        "placement",
			Enabled:     enabled,
			Priority:    10,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000102"),
			Name:        "max-cost-guard",
			Description: "Rejeita cargas cujo custo estimado excede o limite declarado",
			Type:        "cost",
			Enabled:     enabled,
			Priority:    20,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	for i := range demo {
		sharedStore.policies[demo[i].ID] = &demo[i]
	}
}

// HealthCheckHandler returns a simple health check response
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "brlia-waterfall",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ReadyCheckHandler returns a readiness check response
func ReadyCheckHandler(db *pgxpool.Pool, redisClient *redis.Client, eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"service":  "brlia-waterfall",
				"database": "unavailable",
				"error":    err.Error(),
			})
			return
		}

		if _, err := redisClient.Ping(c.Request.Context()).Result(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"service": "brlia-waterfall",
				"redis":   "unavailable",
				"error":   err.Error(),
			})
			return
		}

		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"service":  "brlia-waterfall",
			"database": "available",
			"redis":    "available",
			"engine":   status.Status,
		})
	}
}

// GetEngineStatusHandler exposes the engine status
func GetEngineStatusHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, eng.GetEngineStatus())
	}
}

func workloadToResponse(w *models.Workload) models.WorkloadResponse {
	return models.WorkloadResponse{
		ID:                   w.ID,
		Name:                 w.Name,
		Description:          w.Description,
		OrganizationID:       w.OrganizationID,
		OrganizationName:     "brlia",
		ProjectID:            w.ProjectID,
		ProjectName:          "default",
		UserID:               w.UserID,
		Username:             "system",
		Type:                 w.Type,
		Priority:             w.Priority,
		Status:               w.Status,
		ResourceType:         w.ResourceType,
		ResourceRequirements: w.ResourceRequirements,
		Constraints:          w.Constraints,
		Dependencies:         w.Dependencies,
		Retries:              w.Retries,
		MaxRetries:           w.MaxRetries,
		FailureReason:        w.FailureReason,
		CreatedAt:            w.CreatedAt,
		ScheduledAt:          w.ScheduledAt,
		StartedAt:            w.StartedAt,
		CompletedAt:          w.CompletedAt,
		ExecutionTime:        w.ExecutionTime,
		Deadline:             w.Deadline,
		Timeout:              w.Timeout,
		Metadata:             w.Metadata,
		ClusterID:            w.ClusterID,
		ResourceID:           w.ResourceID,
		Cost:                 w.Cost,
		EstimatedCost:        w.EstimatedCost,
	}
}

// CreateWorkloadHandler creates and submits a new workload
func CreateWorkloadHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.WorkloadCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now().UTC()
		workload := &models.Workload{
			ID:                   uuid.New(),
			Name:                 req.Name,
			Description:          req.Description,
			OrganizationID:       req.OrganizationID,
			ProjectID:            req.ProjectID,
			Type:                 req.Type,
			Priority:             req.Priority,
			Status:               models.WorkloadStatusPending,
			ResourceType:         req.ResourceType,
			ResourceRequirements: req.ResourceRequirements,
			Constraints:          req.Constraints,
			Dependencies:         req.Dependencies,
			MaxRetries:           req.MaxRetries,
			CreatedAt:            now,
			Deadline:             req.Deadline,
			Timeout:              req.Timeout,
			Metadata:             req.Metadata,
		}

		if err := eng.SubmitWorkload(workload); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}

		sharedStore.mu.Lock()
		sharedStore.workloads[workload.ID] = workload
		sharedStore.mu.Unlock()

		c.JSON(http.StatusCreated, workloadToResponse(workload))
	}
}

// ListWorkloadsHandler lists workloads
func ListWorkloadsHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sharedStore.mu.RLock()
		defer sharedStore.mu.RUnlock()

		response := make([]models.WorkloadResponse, 0, len(sharedStore.workloads))
		for _, workload := range sharedStore.workloads {
			response = append(response, workloadToResponse(workload))
		}

		c.JSON(http.StatusOK, gin.H{
			"workloads": response,
			"total":     len(response),
		})
	}
}

// GetWorkloadHandler returns a workload by ID
func GetWorkloadHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workload id"})
			return
		}

		sharedStore.mu.RLock()
		workload, ok := sharedStore.workloads[id]
		sharedStore.mu.RUnlock()

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "workload not found"})
			return
		}

		c.JSON(http.StatusOK, workloadToResponse(workload))
	}
}

// UpdateWorkloadHandler updates mutable workload fields
func UpdateWorkloadHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workload id"})
			return
		}

		var req models.WorkloadUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sharedStore.mu.Lock()
		workload, ok := sharedStore.workloads[id]
		if ok {
			if req.Name != "" {
				workload.Name = req.Name
			}
			if req.Description != "" {
				workload.Description = req.Description
			}
			if req.Priority != 0 {
				workload.Priority = req.Priority
			}
			if req.Status != "" {
				workload.Status = req.Status
			}
			if req.MaxRetries != 0 {
				workload.MaxRetries = req.MaxRetries
			}
			if req.Metadata != nil {
				workload.Metadata = req.Metadata
			}
		}
		sharedStore.mu.Unlock()

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "workload not found"})
			return
		}

		c.JSON(http.StatusOK, workloadToResponse(workload))
	}
}

// DeleteWorkloadHandler removes a workload from the registry
func DeleteWorkloadHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workload id"})
			return
		}

		sharedStore.mu.Lock()
		_, ok := sharedStore.workloads[id]
		delete(sharedStore.workloads, id)
		sharedStore.mu.Unlock()

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "workload not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
