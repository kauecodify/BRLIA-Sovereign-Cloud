package models

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty" db:"organization_id"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	UserID        *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Username      string     `json:"username,omitempty" db:"username"`
	Action        AuditAction `json:"action" db:"action"`
	ResourceType  string     `json:"resource_type" db:"resource_type"`
	ResourceID    uuid.UUID   `json:"resource_id" db:"resource_id"`
	ResourceName  string     `json:"resource_name,omitempty" db:"resource_name"`
	Details       string     `json:"details,omitempty" db:"details"`
	IPAddress     string     `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent     string     `json:"user_agent,omitempty" db:"user_agent"`
	RequestID     string     `json:"request_id,omitempty" db:"request_id"`
	Status        AuditStatus `json:"status" db:"status"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	Metadata      map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// AuditAction represents the type of audit action
type AuditAction string

const (
	AuditActionCreate   AuditAction = "create"
	AuditActionRead     AuditAction = "read"
	AuditActionUpdate   AuditAction = "update"
	AuditActionDelete   AuditAction = "delete"
	AuditActionLogin    AuditAction = "login"
	AuditActionLogout   AuditAction = "logout"
	AuditActionFail     AuditAction = "fail"
	AuditActionSuccess  AuditAction = "success"
	AuditActionWarning  AuditAction = "warning"
)

// AuditStatus represents the status of an audit log
type AuditStatus string

const (
	AuditStatusPending AuditStatus = "pending"
	AuditStatusCompleted AuditStatus = "completed"
	AuditStatusFailed AuditStatus = "failed"
)

// AuditLogCreateRequest represents a request to create an audit log
type AuditLogCreateRequest struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	Username      string     `json:"username,omitempty"`
	Action        AuditAction `json:"action" binding:"required"`
	ResourceType  string     `json:"resource_type" binding:"required"`
	ResourceID    uuid.UUID   `json:"resource_id" binding:"required"`
	ResourceName  string     `json:"resource_name,omitempty"`
	Details       string     `json:"details,omitempty"`
	IPAddress     string     `json:"ip_address,omitempty"`
	UserAgent     string     `json:"user_agent,omitempty"`
	RequestID     string     `json:"request_id,omitempty"`
	Status        AuditStatus `json:"status" binding:"required"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// AuditLogResponse represents an audit log in API responses
type AuditLogResponse struct {
	ID          uuid.UUID    `json:"id"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	OrganizationName *string `json:"organization_name,omitempty"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	ProjectName    *string    `json:"project_name,omitempty"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	Username      string     `json:"username,omitempty"`
	Action        AuditAction `json:"action"`
	ResourceType  string     `json:"resource_type"`
	ResourceID    uuid.UUID   `json:"resource_id"`
	ResourceName  string     `json:"resource_name,omitempty"`
	Details       string     `json:"details,omitempty"`
	IPAddress     string     `json:"ip_address,omitempty"`
	UserAgent     string     `json:"user_agent,omitempty"`
	RequestID     string     `json:"request_id,omitempty"`
	Status        AuditStatus `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// AuditLogListResponse represents a paginated list of audit logs
type AuditLogListResponse struct {
	AuditLogs   []AuditLogResponse `json:"audit_logs"`
	Total       int64              `json:"total"`
	Page        int                `json:"page"`
	PageSize    int                `json:"page_size"`
	TotalPages  int                `json:"total_pages"`
}

// AuditLogFilter represents filters for listing audit logs
type AuditLogFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	Username      string     `json:"username,omitempty"`
	Action        string     `json:"action,omitempty"`
	ResourceType  string     `json:"resource_type,omitempty"`
	ResourceID    *uuid.UUID `json:"resource_id,omitempty"`
	Status        string     `json:"status,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Page         int        `json:"page,omitempty"`
	Limit        int        `json:"limit,omitempty"`
}

// AuditSummary represents a summary of audit logs
type AuditSummary struct {
	TotalActions    int64            `json:"total_actions"`
	ActionsByType   map[string]int64 `json:"actions_by_type"`
	ActionsByUser   map[string]int64 `json:"actions_by_user"`
	ActionsByStatus map[string]int64 `json:"actions_by_status"`
	RecentActions   []AuditLogResponse `json:"recent_actions"`
}
