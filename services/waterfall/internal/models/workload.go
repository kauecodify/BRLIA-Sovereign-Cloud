package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Workload represents a unit of work to be executed
type Workload struct {
	ID                   uuid.UUID            `json:"id" db:"id"`
	Name                 string               `json:"name" db:"name"`
	Description          string               `json:"description,omitempty" db:"description"`
	OrganizationID       uuid.UUID            `json:"organization_id" db:"organization_id"`
	ProjectID            uuid.UUID            `json:"project_id" db:"project_id"`
	UserID               uuid.UUID            `json:"user_id" db:"user_id"`
	Type                 WorkloadType         `json:"type" db:"type"`
	Priority             WorkloadPriority     `json:"priority" db:"priority"`
	Status               WorkloadStatus       `json:"status" db:"status"`
	ResourceType         string               `json:"resource_type" db:"resource_type"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements" db:"resource_requirements"`
	Constraints          WorkloadConstraints  `json:"constraints,omitempty" db:"constraints"`
	Dependencies         []uuid.UUID          `json:"dependencies,omitempty" db:"dependencies"`
	Retries              int                  `json:"retries" db:"retries"`
	MaxRetries           int                  `json:"max_retries,omitempty" db:"max_retries"`
	FailureReason        string               `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt            time.Time            `json:"created_at" db:"created_at"`
	ScheduledAt          *time.Time           `json:"scheduled_at,omitempty" db:"scheduled_at"`
	StartedAt            *time.Time           `json:"started_at,omitempty" db:"started_at"`
	CompletedAt          *time.Time           `json:"completed_at,omitempty" db:"completed_at"`
	ExecutionTime        time.Duration        `json:"execution_time,omitempty" db:"execution_time"`
	Deadline             *time.Time           `json:"deadline,omitempty" db:"deadline"`
	Timeout              time.Duration        `json:"timeout,omitempty" db:"timeout"`
	Metadata             map[string]string    `json:"metadata,omitempty" db:"metadata"`
	ClusterID            *uuid.UUID           `json:"cluster_id,omitempty" db:"cluster_id"`
	ResourceID           *uuid.UUID           `json:"resource_id,omitempty" db:"resource_id"`
	Cost                 float64              `json:"cost,omitempty" db:"cost"`
	EstimatedCost        float64              `json:"estimated_cost,omitempty" db:"estimated_cost"`
}

// WorkloadType represents the type of workload
type WorkloadType string

const (
	WorkloadTypeVM         WorkloadType = "vm"
	WorkloadTypeContainer  WorkloadType = "container"
	WorkloadTypeKubernetes WorkloadType = "kubernetes"
	WorkloadTypeBatch      WorkloadType = "batch"
	WorkloadTypeGPU        WorkloadType = "gpu"
	WorkloadTypeHPC        WorkloadType = "hpc"
	WorkloadTypeServerless WorkloadType = "serverless"
	WorkloadTypeAI         WorkloadType = "ai"
	WorkloadTypeQuantum    WorkloadType = "quantum"
)

// WorkloadPriority represents the priority of a workload
type WorkloadPriority int

const (
	WorkloadPriorityLowest  WorkloadPriority = 0
	WorkloadPriorityLow     WorkloadPriority = 1
	WorkloadPriorityNormal  WorkloadPriority = 2
	WorkloadPriorityHigh    WorkloadPriority = 3
	WorkloadPriorityHighest WorkloadPriority = 4
)

// WorkloadStatus represents the current status of a workload
type WorkloadStatus string

const (
	WorkloadStatusPending    WorkloadStatus = "pending"
	WorkloadStatusScheduled  WorkloadStatus = "scheduled"
	WorkloadStatusProcessing WorkloadStatus = "processing"
	WorkloadStatusCompleted  WorkloadStatus = "completed"
	WorkloadStatusFailed     WorkloadStatus = "failed"
	WorkloadStatusCancelled  WorkloadStatus = "cancelled"
	WorkloadStatusPaused     WorkloadStatus = "paused"
)

// ResourceRequirements represents the resource requirements for a workload
type ResourceRequirements struct {
	CPU       float64 `json:"cpu"`
	MemoryGB  float64 `json:"memory_gb"`
	StorageGB float64 `json:"storage_gb"`
	GPUCount  int     `json:"gpu_count"`
	GPUType   string  `json:"gpu_type,omitempty"`
	Bandwidth float64 `json:"bandwidth_mbps,omitempty"`
	MinNodes  int     `json:"min_nodes,omitempty"`
	MaxNodes  int     `json:"max_nodes,omitempty"`
}

// WorkloadConstraints represents constraints for workload scheduling
type WorkloadConstraints struct {
	RequiredRegion  *string          `json:"required_region,omitempty"`
	RequiredZone    *string          `json:"required_zone,omitempty"`
	RequiredCluster *uuid.UUID       `json:"required_cluster,omitempty"`
	AvoidRegions    []string         `json:"avoid_regions,omitempty"`
	AvoidZones      []string         `json:"avoid_zones,omitempty"`
	AvoidClusters   []uuid.UUID      `json:"avoid_clusters,omitempty"`
	RequiredTags    []string         `json:"required_tags,omitempty"`
	AvoidTags       []string         `json:"avoid_tags,omitempty"`
	MaxCost         *float64         `json:"max_cost,omitempty"`
	SLARequirements *SLARequirements `json:"sla_requirements,omitempty"`
}

// SLARequirements represents SLA requirements for a workload
type SLARequirements struct {
	Availability float64 `json:"availability,omitempty"`
	LatencyMS    int64   `json:"latency_ms,omitempty"`
	Throughput   float64 `json:"throughput,omitempty"`
}

// WorkloadCreateRequest represents a request to create a workload
type WorkloadCreateRequest struct {
	Name                 string               `json:"name" binding:"required,min=1,max=100"`
	Description          string               `json:"description,omitempty" binding:"max=500"`
	OrganizationID       uuid.UUID            `json:"organization_id" binding:"required"`
	ProjectID            uuid.UUID            `json:"project_id" binding:"required"`
	Type                 WorkloadType         `json:"type" binding:"required,oneof=vm container kubernetes batch gpu hpc serverless ai quantum"`
	Priority             WorkloadPriority     `json:"priority,omitempty" binding:"min=0,max=4"`
	ResourceType         string               `json:"resource_type" binding:"required"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements" binding:"required"`
	Constraints          WorkloadConstraints  `json:"constraints,omitempty"`
	Dependencies         []uuid.UUID          `json:"dependencies,omitempty"`
	MaxRetries           int                  `json:"max_retries,omitempty" binding:"min=0,max=10"`
	Deadline             *time.Time           `json:"deadline,omitempty"`
	Timeout              time.Duration        `json:"timeout,omitempty"`
	Metadata             map[string]string    `json:"metadata,omitempty"`
}

// WorkloadUpdateRequest represents a request to update a workload
type WorkloadUpdateRequest struct {
	Name                 string               `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description          string               `json:"description,omitempty" binding:"omitempty,max=500"`
	Priority             WorkloadPriority     `json:"priority,omitempty" binding:"omitempty,min=0,max=4"`
	Status               WorkloadStatus       `json:"status,omitempty" binding:"omitempty,oneof=pending scheduled processing completed failed cancelled paused"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements,omitempty"`
	Constraints          WorkloadConstraints  `json:"constraints,omitempty"`
	MaxRetries           int                  `json:"max_retries,omitempty" binding:"omitempty,min=0,max=10"`
	Deadline             *time.Time           `json:"deadline,omitempty"`
	Timeout              time.Duration        `json:"timeout,omitempty"`
	Metadata             map[string]string    `json:"metadata,omitempty"`
}

// WorkloadResponse represents a workload in API responses
type WorkloadResponse struct {
	ID                   uuid.UUID            `json:"id"`
	Name                 string               `json:"name"`
	Description          string               `json:"description,omitempty"`
	OrganizationID       uuid.UUID            `json:"organization_id"`
	OrganizationName     string               `json:"organization_name"`
	ProjectID            uuid.UUID            `json:"project_id"`
	ProjectName          string               `json:"project_name"`
	UserID               uuid.UUID            `json:"user_id"`
	Username             string               `json:"username"`
	Type                 WorkloadType         `json:"type"`
	Priority             WorkloadPriority     `json:"priority"`
	Status               WorkloadStatus       `json:"status"`
	ResourceType         string               `json:"resource_type"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements"`
	Constraints          WorkloadConstraints  `json:"constraints,omitempty"`
	Dependencies         []uuid.UUID          `json:"dependencies,omitempty"`
	Retries              int                  `json:"retries"`
	MaxRetries           int                  `json:"max_retries,omitempty"`
	FailureReason        string               `json:"failure_reason,omitempty"`
	CreatedAt            time.Time            `json:"created_at"`
	ScheduledAt          *time.Time           `json:"scheduled_at,omitempty"`
	StartedAt            *time.Time           `json:"started_at,omitempty"`
	CompletedAt          *time.Time           `json:"completed_at,omitempty"`
	ExecutionTime        time.Duration        `json:"execution_time,omitempty"`
	Deadline             *time.Time           `json:"deadline,omitempty"`
	Timeout              time.Duration        `json:"timeout,omitempty"`
	Metadata             map[string]string    `json:"metadata,omitempty"`
	ClusterID            *uuid.UUID           `json:"cluster_id,omitempty"`
	ClusterName          *string              `json:"cluster_name,omitempty"`
	ResourceID           *uuid.UUID           `json:"resource_id,omitempty"`
	ResourceName         *string              `json:"resource_name,omitempty"`
	Cost                 float64              `json:"cost,omitempty"`
	EstimatedCost        float64              `json:"estimated_cost,omitempty"`
}

// WorkloadListResponse represents a paginated list of workloads
type WorkloadListResponse struct {
	Workloads  []WorkloadResponse `json:"workloads"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// WorkloadFilter represents filters for listing workloads
type WorkloadFilter struct {
	OrganizationID uuid.UUID        `json:"organization_id,omitempty"`
	ProjectID      uuid.UUID        `json:"project_id,omitempty"`
	UserID         uuid.UUID        `json:"user_id,omitempty"`
	Name           string           `json:"name,omitempty"`
	Type           WorkloadType     `json:"type,omitempty"`
	Priority       WorkloadPriority `json:"priority,omitempty"`
	Status         WorkloadStatus   `json:"status,omitempty"`
	ResourceType   string           `json:"resource_type,omitempty"`
	ClusterID      uuid.UUID        `json:"cluster_id,omitempty"`
	StartTime      *time.Time       `json:"start_time,omitempty"`
	EndTime        *time.Time       `json:"end_time,omitempty"`
	Page           int              `json:"page,omitempty"`
	Limit          int              `json:"limit,omitempty"`
}

// Validate validates the workload
func (w *Workload) Validate() error {
	if w.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(w.Name) > 100 {
		return fmt.Errorf("name must be 100 characters or less")
	}
	if w.OrganizationID == uuid.Nil {
		return fmt.Errorf("organization_id is required")
	}
	if w.ProjectID == uuid.Nil {
		return fmt.Errorf("project_id is required")
	}
	if w.Type == "" {
		return fmt.Errorf("type is required")
	}
	if w.ResourceType == "" {
		return fmt.Errorf("resource_type is required")
	}
	if w.ResourceRequirements.CPU <= 0 {
		return fmt.Errorf("CPU requirement must be positive")
	}
	if w.Priority < WorkloadPriorityLowest || w.Priority > WorkloadPriorityHighest {
		return fmt.Errorf("priority must be between 0 and 4")
	}
	return nil
}

// WorkloadMetrics represents metrics for a workload
type WorkloadMetrics struct {
	WorkloadID   uuid.UUID      `json:"workload_id"`
	CPUUsage     float64        `json:"cpu_usage_percent"`
	MemoryUsage  float64        `json:"memory_usage_percent"`
	StorageUsage float64        `json:"storage_usage_percent"`
	GPUUsage     float64        `json:"gpu_usage_percent,omitempty"`
	NetworkIn    float64        `json:"network_in_mbps"`
	NetworkOut   float64        `json:"network_out_mbps"`
	Status       WorkloadStatus `json:"status"`
	LastUpdated  time.Time      `json:"last_updated"`
}
