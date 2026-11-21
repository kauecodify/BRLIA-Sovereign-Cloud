package handlers

import (
	"net/http"
	"time"

	"github.com/brlia/waterfall/internal/engine"
	"github.com/brlia/waterfall/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListResourcesHandler lists known infrastructure resources
func ListResourcesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"resources": demoResources(),
			"total":     len(demoResources()),
		})
	}
}

// GetResourceHandler returns a resource by ID
func GetResourceHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
			return
		}

		for _, resource := range demoResources() {
			if resource.ID == id {
				c.JSON(http.StatusOK, resource)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	}
}

// GetResourceByTypeHandler returns resources by type filter
func GetResourcesByTypeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceType := c.Param("type")
		filtered := make([]models.Resource, 0)
		for _, resource := range demoResources() {
			if resource.Type == resourceType {
				filtered = append(filtered, resource)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"resources": filtered,
			"total":     len(filtered),
			"type":      resourceType,
		})
	}
}

func demoResources() []models.Resource {
	return []models.Resource{
		{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:         "brlia-node-br-sul-1a",
			Type:         "kvm",
			Region:       "br-sul-1",
			Zone:         "a",
			Status:       "available",
			TotalCPU:     64,
			TotalMemory:  256,
			TotalStorage: 4096,
			TotalGPU:     4,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:         "brlia-node-br-sul-1b",
			Type:         "kvm",
			Region:       "br-sul-1",
			Zone:         "b",
			Status:       "available",
			TotalCPU:     64,
			TotalMemory:  256,
			TotalStorage: 4096,
			TotalGPU:     4,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}
}

// GetCapacityHandler returns the current capacity summary
func GetCapacityHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{
			"resources":      status.Resources,
			"workloads":      status.ActiveWorkloads,
			"pending":        status.PendingWorkloads,
			"last_heartbeat": status.LastHeartbeat,
		})
	}
}

// GetCapacityForecastHandler returns a capacity forecast
func GetCapacityForecastHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{
			"forecast": []gin.H{
				{
					"horizon":       "7d",
					"projected_cpu": status.Resources.UsedCPU,
					"available_cpu": status.Resources.AvailableCPU,
				},
				{
					"horizon":       "30d",
					"projected_cpu": status.Resources.UsedCPU,
					"available_cpu": status.Resources.AvailableCPU,
				},
			},
			"generated_at": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ListQueuesHandler lists scheduling queues
func ListQueuesHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"queues": []gin.H{
				{"name": "default", "description": "Fila padrao para cargas de uso geral", "workers": 10},
				{"name": "priority", "description": "Fila de prioridade para cargas criticas", "workers": 5},
			},
		})
	}
}

// GetQueueHandler returns a queue by name
func GetQueueHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		queues := map[string]gin.H{
			"default":  {"name": "default", "description": "Fila padrao para cargas de uso geral", "workers": 10},
			"priority": {"name": "priority", "description": "Fila de prioridade para cargas criticas", "workers": 5},
		}
		queue, ok := queues[name]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "queue not found"})
			return
		}
		c.JSON(http.StatusOK, queue)
	}
}

// ScheduleWorkloadHandler submits a workload for immediate scheduling
func ScheduleWorkloadHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
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
			OrganizationID:       req.OrganizationID,
			ProjectID:            req.ProjectID,
			Type:                 req.Type,
			Priority:             req.Priority,
			Status:               models.WorkloadStatusScheduled,
			ResourceType:         req.ResourceType,
			ResourceRequirements: req.ResourceRequirements,
			MaxRetries:           req.MaxRetries,
			CreatedAt:            now,
		}

		if err := eng.SubmitWorkload(workload); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}

		sharedStore.mu.Lock()
		sharedStore.workloads[workload.ID] = workload
		sharedStore.mu.Unlock()

		c.JSON(http.StatusAccepted, workloadToResponse(workload))
	}
}

// GetScheduleHandler returns the current schedule view
func GetScheduleHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{
			"active_workloads":  status.ActiveWorkloads,
			"pending_workloads": status.PendingWorkloads,
			"recent_actions":    status.RecentActions,
		})
	}
}

// ListPoliciesHandler lists scheduling and governance policies
func ListPoliciesHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sharedStore.mu.RLock()
		defer sharedStore.mu.RUnlock()

		policies := make([]models.Policy, 0, len(sharedStore.policies))
		for _, policy := range sharedStore.policies {
			policies = append(policies, *policy)
		}
		c.JSON(http.StatusOK, gin.H{"policies": policies, "total": len(policies)})
	}
}

// CreatePolicyHandler registers a new policy
func CreatePolicyHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PolicyCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		enabled := true
		if req.Enabled != nil {
			enabled = *req.Enabled
		}
		now := time.Now().UTC()
		policy := models.Policy{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Enabled:     enabled,
			Priority:    req.Priority,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		sharedStore.mu.Lock()
		sharedStore.policies[policy.ID] = &policy
		sharedStore.mu.Unlock()

		c.JSON(http.StatusCreated, policy)
	}
}

// GetPolicyHandler returns a policy by ID
func GetPolicyHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
			return
		}

		sharedStore.mu.RLock()
		policy, ok := sharedStore.policies[id]
		sharedStore.mu.RUnlock()

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}

		c.JSON(http.StatusOK, policy)
	}
}

// UpdatePolicyHandler updates a policy
func UpdatePolicyHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
			return
		}

		var req models.PolicyUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sharedStore.mu.Lock()
		policy, ok := sharedStore.policies[id]
		if ok {
			if req.Name != "" {
				policy.Name = req.Name
			}
			if req.Description != "" {
				policy.Description = req.Description
			}
			if req.Enabled != nil {
				policy.Enabled = *req.Enabled
			}
			if req.Priority != nil {
				policy.Priority = *req.Priority
			}
			policy.UpdatedAt = time.Now().UTC()
		}
		sharedStore.mu.Unlock()

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}

		c.JSON(http.StatusOK, policy)
	}
}

// DeletePolicyHandler removes a policy
func DeletePolicyHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
			return
		}

		sharedStore.mu.Lock()
		_, ok := sharedStore.policies[id]
		delete(sharedStore.policies, id)
		sharedStore.mu.Unlock()

		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetMetricsHandler returns engine-level metrics
func GetMetricsHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{
			"metrics":        status.Metrics,
			"worker_count":   status.WorkerCount,
			"last_heartbeat": status.LastHeartbeat,
		})
	}
}

// GetWorkloadMetricsHandler returns per-workload metrics
func GetWorkloadMetricsHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sharedStore.mu.RLock()
		defer sharedStore.mu.RUnlock()

		metrics := make([]gin.H, 0, len(sharedStore.workloads))
		for _, workload := range sharedStore.workloads {
			metrics = append(metrics, gin.H{
				"workload_id": workload.ID,
				"status":      workload.Status,
				"retries":     workload.Retries,
				"cost":        workload.Cost,
			})
		}
		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

// GetResourceMetricsHandler returns per-resource metrics
func GetResourceMetricsHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{
			"resources": status.Resources,
			"clusters":  status.Clusters,
		})
	}
}

// GetCostOptimizationHandler returns cost optimization suggestions
func GetCostOptimizationHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"suggestions": []gin.H{
				{
					"type":                  "consolidation",
					"description":           "Consolidar cargas com baixa utilizacao de CPU no cluster br-sul-1",
					"estimated_savings_brl": 1250.00,
				},
			},
			"currency": "BRL",
		})
	}
}

// AnalyzeCostsHandler analyzes costs for a period
func AnalyzeCostsHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sharedStore.mu.RLock()
		totalCost := 0.0
		for _, workload := range sharedStore.workloads {
			totalCost += workload.EstimatedCost
		}
		sharedStore.mu.RUnlock()

		c.JSON(http.StatusOK, gin.H{
			"total_estimated_cost_brl": totalCost,
			"currency":                 "BRL",
			"analyzed_at":              time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ListFailuresHandler lists detected failures
func ListFailuresHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sharedStore.mu.RLock()
		defer sharedStore.mu.RUnlock()

		failures := make([]gin.H, 0)
		for _, workload := range sharedStore.workloads {
			if workload.Status == models.WorkloadStatusFailed {
				failures = append(failures, gin.H{
					"id":             workload.ID,
					"workload":       workload.Name,
					"failure_reason": workload.FailureReason,
					"retries":        workload.Retries,
					"detected_at":    workload.CreatedAt,
				})
			}
		}
		c.JSON(http.StatusOK, gin.H{"failures": failures, "total": len(failures)})
	}
}

// GetFailureHandler returns failure details for a workload
func GetFailureHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid failure id"})
			return
		}

		sharedStore.mu.RLock()
		workload, ok := sharedStore.workloads[id]
		sharedStore.mu.RUnlock()

		if !ok || workload.Status != models.WorkloadStatusFailed {
			c.JSON(http.StatusNotFound, gin.H{"error": "failure not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":             workload.ID,
			"workload":       workload.Name,
			"failure_reason": workload.FailureReason,
			"retries":        workload.Retries,
		})
	}
}

// TriggerRecoveryHandler triggers a recovery procedure
func TriggerRecoveryHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		sharedStore.mu.Lock()
		recovered := make([]uuid.UUID, 0)
		for id, workload := range sharedStore.workloads {
			if workload.Status == models.WorkloadStatusFailed && workload.Retries < workload.MaxRetries {
				workload.Status = models.WorkloadStatusPending
				workload.FailureReason = ""
				recovered = append(recovered, id)
			}
		}
		sharedStore.mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"triggered": true,
			"requeued":  recovered,
			"count":     len(recovered),
		})
	}
}

// ListClustersHandler lists known clusters
func ListClustersHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := eng.GetEngineStatus()
		c.JSON(http.StatusOK, gin.H{"clusters": status.Clusters, "total": len(status.Clusters)})
	}
}

// GetClusterHandler returns a cluster by ID
func GetClusterHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
			return
		}

		status := eng.GetEngineStatus()
		for _, cluster := range status.Clusters {
			if cluster.ID == id {
				c.JSON(http.StatusOK, cluster)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
	}
}

// AssignWorkloadToClusterHandler pins a workload to a cluster
func AssignWorkloadToClusterHandler(eng *engine.WaterfallEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
			return
		}

		var req models.WorkloadCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now().UTC()
		workload := &models.Workload{
			ID:                   uuid.New(),
			Name:                 req.Name,
			OrganizationID:       req.OrganizationID,
			ProjectID:            req.ProjectID,
			Type:                 req.Type,
			Status:               models.WorkloadStatusScheduled,
			ResourceType:         req.ResourceType,
			ResourceRequirements: req.ResourceRequirements,
			CreatedAt:            now,
			ClusterID:            &clusterID,
		}

		if err := eng.SubmitWorkload(workload); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}

		sharedStore.mu.Lock()
		sharedStore.workloads[workload.ID] = workload
		sharedStore.mu.Unlock()

		c.JSON(http.StatusAccepted, workloadToResponse(workload))
	}
}
