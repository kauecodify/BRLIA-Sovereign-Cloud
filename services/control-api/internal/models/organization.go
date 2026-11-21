package models

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents a top-level container for projects and users
type Organization struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Description string             `json:"description,omitempty" db:"description"`
	Status      OrganizationStatus `json:"status" db:"status"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty" db:"deleted_at"`
	Metadata    map[string]string  `json:"metadata,omitempty" db:"metadata"`
}

// OrganizationStatus represents the current state of an organization
type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusInactive  OrganizationStatus = "inactive"
	OrganizationStatusPending   OrganizationStatus = "pending"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
)

// OrganizationCreateRequest represents a request to create an organization
type OrganizationCreateRequest struct {
	Name        string            `json:"name" binding:"required,min=1,max=100"`
	Description string            `json:"description,omitempty" binding:"max=500"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// OrganizationUpdateRequest represents a request to update an organization
type OrganizationUpdateRequest struct {
	Name        string             `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description string             `json:"description,omitempty" binding:"omitempty,max=500"`
	Status      OrganizationStatus `json:"status,omitempty" binding:"omitempty,oneof=active inactive pending suspended"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
}

// OrganizationResponse represents an organization in API responses
type OrganizationResponse struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description,omitempty"`
	Status       OrganizationStatus `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Metadata     map[string]string  `json:"metadata,omitempty"`
	ProjectCount int                `json:"project_count"`
	UserCount    int                `json:"user_count"`
}

// OrganizationListResponse represents a paginated list of organizations
type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	TotalPages    int                    `json:"total_pages"`
}

// OrganizationFilter represents filters for listing organizations
type OrganizationFilter struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}
