package models

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a set of permissions
type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Permissions []string `json:"permissions" db:"permissions"`
	IsSystem    bool      `json:"is_system" db:"is_system"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// RoleCreateRequest represents a request to create a role
type RoleCreateRequest struct {
	Name        string   `json:"name" binding:"required,min=1,max=50"`
	Description string   `json:"description,omitempty" binding:"max=500"`
	Permissions []string `json:"permissions" binding:"required"`
}

// RoleUpdateRequest represents a request to update a role
type RoleUpdateRequest struct {
	Name        string   `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Description string   `json:"description,omitempty" binding:"omitempty,max=500"`
	Permissions []string `json:"permissions,omitempty"`
}

// RoleResponse represents a role in API responses
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Permissions []string `json:"permissions"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserCount   int       `json:"user_count"`
}

// RoleListResponse represents a paginated list of roles
type RoleListResponse struct {
	Roles     []RoleResponse `json:"roles"`
	Total     int64         `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

// RoleFilter represents filters for listing roles
type RoleFilter struct {
	Name     string `json:"name,omitempty"`
	IsSystem bool   `json:"is_system,omitempty"`
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// Permission represents a single permission
type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Effect   string `json:"effect"`
}

// Common permissions
const (
	PermissionRead   = "read"
	PermissionWrite  = "write"
	PermissionDelete = "delete"
	PermissionAdmin  = "admin"
	PermissionAll    = "*"
)

// Common resource types
const (
	ResourceOrganization = "organization"
	ResourceProject      = "project"
	ResourceUser         = "user"
	ResourceRole         = "role"
	ResourceQuota        = "quota"
	ResourceRegion       = "region"
	ResourceCompute      = "compute"
	ResourceStorage      = "storage"
	ResourceNetwork      = "network"
	ResourceBilling      = "billing"
	ResourceAudit        = "audit"
)

// System roles
const (
	RoleSuperAdmin = "super-admin"
	RoleAdmin       = "admin"
	RoleEditor      = "editor"
	RoleViewer      = "viewer"
	RoleUser        = "user"
)
