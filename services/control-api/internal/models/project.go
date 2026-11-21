package models

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a container for resources within an organization
type Project struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	OrganizationID uuid.UUID    `json:"organization_id" db:"organization_id"`
	Name           string       `json:"name" db:"name"`
	Description    string       `json:"description,omitempty" db:"description"`
	Status         ProjectStatus `json:"status" db:"status"`
	RegionID       uuid.UUID    `json:"region_id" db:"region_id"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
	Metadata       map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// ProjectStatus represents the current state of a project
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusInactive ProjectStatus = "inactive"
	ProjectStatusPending  ProjectStatus = "pending"
	ProjectStatusSuspended ProjectStatus = "suspended"
	ProjectStatusArchived ProjectStatus = "archived"
)

// ProjectCreateRequest represents a request to create a project
type ProjectCreateRequest struct {
	OrganizationID uuid.UUID    `json:"organization_id" binding:"required"`
	Name           string       `json:"name" binding:"required,min=1,max=100"`
	Description    string       `json:"description,omitempty" binding:"max=500"`
	RegionID       uuid.UUID    `json:"region_id" binding:"required"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ProjectUpdateRequest represents a request to update a project
type ProjectUpdateRequest struct {
	Name        string       `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description string       `json:"description,omitempty" binding:"omitempty,max=500"`
	Status      ProjectStatus `json:"status,omitempty" binding:"omitempty,oneof=active inactive pending suspended archived"`
	RegionID    uuid.UUID    `json:"region_id,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ProjectResponse represents a project in API responses
type ProjectResponse struct {
	ID             uuid.UUID    `json:"id"`
	OrganizationID uuid.UUID    `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	Name           string       `json:"name"`
	Description    string       `json:"description,omitempty"`
	Status         ProjectStatus `json:"status"`
	RegionID       uuid.UUID    `json:"region_id"`
	RegionName     string       `json:"region_name"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	ResourceCount  int          `json:"resource_count"`
	UserCount      int          `json:"user_count"`
}

// ProjectListResponse represents a paginated list of projects
type ProjectListResponse struct {
	Projects    []ProjectResponse `json:"projects"`
	Total       int64             `json:"total"`
	Page        int               `json:"page"`
	PageSize    int               `json:"page_size"`
	TotalPages  int               `json:"total_pages"`
}

// ProjectFilter represents filters for listing projects
type ProjectFilter struct {
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Status         string    `json:"status,omitempty"`
	RegionID       uuid.UUID `json:"region_id,omitempty"`
	Page           int       `json:"page,omitempty"`
	Limit          int       `json:"limit,omitempty"`
}
