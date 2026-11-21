package models

import (
	"time"

	"github.com/google/uuid"
)

// Resource represents an infrastructure resource managed by the engine
type Resource struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Type         string     `json:"type" db:"type"`
	Region       string     `json:"region" db:"region"`
	Zone         string     `json:"zone" db:"zone"`
	Status       string     `json:"status" db:"status"`
	ClusterID    *uuid.UUID `json:"cluster_id,omitempty" db:"cluster_id"`
	TotalCPU     float64    `json:"total_cpu" db:"total_cpu"`
	UsedCPU      float64    `json:"used_cpu" db:"used_cpu"`
	TotalMemory  float64    `json:"total_memory_gb" db:"total_memory_gb"`
	UsedMemory   float64    `json:"used_memory_gb" db:"used_memory_gb"`
	TotalStorage float64    `json:"total_storage_gb" db:"total_storage_gb"`
	UsedStorage  float64    `json:"used_storage_gb" db:"used_storage_gb"`
	TotalGPU     int        `json:"total_gpu" db:"total_gpu"`
	UsedGPU      int        `json:"used_gpu" db:"used_gpu"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// ResourceListResponse represents a paginated list of resources
type ResourceListResponse struct {
	Resources []Resource `json:"resources"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
}

// QueueInfo represents a workload queue
type QueueInfo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Depth         int    `json:"depth"`
	Workers       int    `json:"workers"`
	WorkloadTypes string `json:"workload_types"`
}

// Policy represents a scheduling or governance policy
type Policy struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`
	Enabled     bool      `json:"enabled" db:"enabled"`
	Priority    int       `json:"priority" db:"priority"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PolicyCreateRequest represents a request to create a policy
type PolicyCreateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description,omitempty" binding:"max=500"`
	Type        string `json:"type" binding:"required,oneof=scheduling capacity cost security placement"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Priority    int    `json:"priority,omitempty" binding:"min=0,max=100"`
}

// PolicyUpdateRequest represents a request to update a policy
type PolicyUpdateRequest struct {
	Name        string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description string `json:"description,omitempty" binding:"omitempty,max=500"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Priority    *int   `json:"priority,omitempty" binding:"omitempty,min=0,max=100"`
}
