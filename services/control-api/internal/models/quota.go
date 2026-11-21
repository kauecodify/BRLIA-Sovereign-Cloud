package models

import (
	"time"

	"github.com/google/uuid"
)

// Quota represents resource limits for an organization or project
type Quota struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	OrganizationID *uuid.UUID   `json:"organization_id,omitempty" db:"organization_id"`
	ProjectID      *uuid.UUID   `json:"project_id,omitempty" db:"project_id"`
	ResourceType   string       `json:"resource_type" db:"resource_type"`
	Limit          int64        `json:"limit" db:"limit"`
	Used           int64        `json:"used" db:"used"`
	Unit           string       `json:"unit" db:"unit"`
	Scope          QuotaScope   `json:"scope" db:"scope"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
	Metadata       map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// QuotaScope represents the scope of a quota
type QuotaScope string

const (
	QuotaScopeOrganization QuotaScope = "organization"
	QuotaScopeProject      QuotaScope = "project"
	QuotaScopeGlobal       QuotaScope = "global"
)

// Resource types for quotas
const (
	ResourceTypeVM         = "vm"
	ResourceTypeCPU        = "cpu"
	ResourceTypeMemory     = "memory"
	ResourceTypeStorage    = "storage"
	ResourceTypeGPU        = "gpu"
	ResourceTypeIP         = "ip"
	ResourceTypeBandwidth  = "bandwidth"
	ResourceTypeContainers = "containers"
	ResourceTypeKubernetes = "kubernetes"
	ResourceTypeDatabases  = "databases"
	ResourceTypeUsers      = "users"
)

// QuotaCreateRequest represents a request to create a quota
type QuotaCreateRequest struct {
	OrganizationID *uuid.UUID   `json:"organization_id,omitempty"`
	ProjectID      *uuid.UUID   `json:"project_id,omitempty"`
	ResourceType   string       `json:"resource_type" binding:"required"`
	Limit          int64        `json:"limit" binding:"required,min=0"`
	Unit           string       `json:"unit" binding:"required"`
	Scope          QuotaScope   `json:"scope" binding:"required,oneof=organization project global"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// QuotaUpdateRequest represents a request to update a quota
type QuotaUpdateRequest struct {
	Limit    int64        `json:"limit,omitempty" binding:"omitempty,min=0"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// QuotaResponse represents a quota in API responses
type QuotaResponse struct {
	ID             uuid.UUID    `json:"id"`
	OrganizationID *uuid.UUID   `json:"organization_id,omitempty"`
	OrganizationName *string   `json:"organization_name,omitempty"`
	ProjectID      *uuid.UUID   `json:"project_id,omitempty"`
	ProjectName    *string      `json:"project_name,omitempty"`
	ResourceType   string       `json:"resource_type"`
	Limit          int64        `json:"limit"`
	Used           int64        `json:"used"`
	Available      int64        `json:"available"`
	Unit           string       `json:"unit"`
	Scope          QuotaScope   `json:"scope"`
	UsagePercent   float64      `json:"usage_percent"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// QuotaListResponse represents a paginated list of quotas
type QuotaListResponse struct {
	Quotas     []QuotaResponse `json:"quotas"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// QuotaFilter represents filters for listing quotas
type QuotaFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	ResourceType   string     `json:"resource_type,omitempty"`
	Scope          string     `json:"scope,omitempty"`
	Page           int        `json:"page,omitempty"`
	Limit          int        `json:"limit,omitempty"`
}

// QuotaUsage represents the current usage of a quota
type QuotaUsage struct {
	ResourceType string `json:"resource_type"`
	Used         int64  `json:"used"`
	Limit        int64  `json:"limit"`
	Unit         string `json:"unit"`
	Percent      float64 `json:"percent"`
}
