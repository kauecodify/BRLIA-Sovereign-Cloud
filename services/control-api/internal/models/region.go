package models

import (
	"time"

	"github.com/google/uuid"
)

// Region represents a geographic region with availability zones
type Region struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Location    string    `json:"location" db:"location"`
	Country     string    `json:"country" db:"country"`
	Status      RegionStatus `json:"status" db:"status"`
	IsDefault   bool      `json:"is_default" db:"is_default"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Metadata    map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// RegionStatus represents the current state of a region
type RegionStatus string

const (
	RegionStatusActive   RegionStatus = "active"
	RegionStatusInactive RegionStatus = "inactive"
	RegionStatusMaintenance RegionStatus = "maintenance"
)

// AvailabilityZone represents an availability zone within a region
type AvailabilityZone struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RegionID    uuid.UUID `json:"region_id" db:"region_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Status      ZoneStatus `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// ZoneStatus represents the current state of an availability zone
type ZoneStatus string

const (
	ZoneStatusActive   ZoneStatus = "active"
	ZoneStatusInactive ZoneStatus = "inactive"
	ZoneStatusMaintenance ZoneStatus = "maintenance"
)

// RegionCreateRequest represents a request to create a region
type RegionCreateRequest struct {
	Name        string          `json:"name" binding:"required,min=1,max=50"`
	Description string          `json:"description,omitempty" binding:"max=500"`
	Location    string          `json:"location" binding:"required"`
	Country     string          `json:"country" binding:"required"`
	IsDefault   bool            `json:"is_default,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Zones       []ZoneCreateRequest `json:"zones,omitempty"`
}

// ZoneCreateRequest represents a request to create an availability zone
type ZoneCreateRequest struct {
	Name        string          `json:"name" binding:"required,min=1,max=50"`
	Description string          `json:"description,omitempty" binding:"max=500"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// RegionUpdateRequest represents a request to update a region
type RegionUpdateRequest struct {
	Name        string          `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Description string          `json:"description,omitempty" binding:"omitempty,max=500"`
	Location    string          `json:"location,omitempty"`
	Country     string          `json:"country,omitempty"`
	Status      RegionStatus    `json:"status,omitempty" binding:"omitempty,oneof=active inactive maintenance"`
	IsDefault   bool            `json:"is_default,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// RegionResponse represents a region in API responses
type RegionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location"`
	Country     string    `json:"country"`
	Status      RegionStatus `json:"status"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ZoneCount   int       `json:"zone_count"`
	ClusterCount int      `json:"cluster_count"`
	Zones       []ZoneResponse `json:"zones,omitempty"`
}

// ZoneResponse represents an availability zone in API responses
type ZoneResponse struct {
	ID          uuid.UUID `json:"id"`
	RegionID    uuid.UUID `json:"region_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      ZoneStatus `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ClusterCount int      `json:"cluster_count"`
}

// RegionListResponse represents a paginated list of regions
type RegionListResponse struct {
	Regions    []RegionResponse `json:"regions"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// RegionFilter represents filters for listing regions
type RegionFilter struct {
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
	Status  string `json:"status,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

// Cluster represents a compute cluster in a region/zone
type Cluster struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RegionID    uuid.UUID `json:"region_id" db:"region_id"`
	ZoneID      uuid.UUID `json:"zone_id" db:"zone_id"`
	Name        string    `json:"name" db:"name"`
	Type        ClusterType `json:"type" db:"type"`
	Status      ClusterStatus `json:"status" db:"status"`
	TotalCPU    int64     `json:"total_cpu" db:"total_cpu"`
	TotalMemory int64     `json:"total_memory" db:"total_memory"`
	TotalStorage int64   `json:"total_storage" db:"total_storage"`
	UsedCPU     int64     `json:"used_cpu" db:"used_cpu"`
	UsedMemory  int64     `json:"used_memory" db:"used_memory"`
	UsedStorage int64     `json:"used_storage" db:"used_storage"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// ClusterType represents the type of a cluster
type ClusterType string

const (
	ClusterTypeKubernetes ClusterType = "kubernetes"
	ClusterTypeKVM        ClusterType = "kvm"
	ClusterTypeBareMetal  ClusterType = "bare-metal"
	ClusterTypeHPC        ClusterType = "hpc"
	ClusterTypeGPU        ClusterType = "gpu"
)

// ClusterStatus represents the current state of a cluster
type ClusterStatus string

const (
	ClusterStatusActive   ClusterStatus = "active"
	ClusterStatusInactive ClusterStatus = "inactive"
	ClusterStatusDraining ClusterStatus = "draining"
	ClusterStatusError    ClusterStatus = "error"
)
