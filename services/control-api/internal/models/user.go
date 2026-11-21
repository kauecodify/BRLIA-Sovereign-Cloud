package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID             uuid.UUID         `json:"id" db:"id"`
	OrganizationID uuid.UUID         `json:"organization_id" db:"organization_id"`
	Email          string            `json:"email" db:"email"`
	Username       string            `json:"username" db:"username"`
	DisplayName    string            `json:"display_name" db:"display_name"`
	PasswordHash   string            `json:"-" db:"password_hash"`
	Status         UserStatus        `json:"status" db:"status"`
	MFAEnabled     bool              `json:"mfa_enabled" db:"mfa_enabled"`
	MFASecret      string            `json:"-" db:"mfa_secret"`
	LastLogin      *time.Time        `json:"last_login,omitempty" db:"last_login"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`
	Metadata       map[string]string `json:"metadata,omitempty" db:"metadata"`
}

// UserStatus represents the current state of a user
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusPending   UserStatus = "pending"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusLocked    UserStatus = "locked"
)

// UserCreateRequest represents a request to create a user
type UserCreateRequest struct {
	OrganizationID uuid.UUID         `json:"organization_id" binding:"required"`
	Email          string            `json:"email" binding:"required,email"`
	Username       string            `json:"username" binding:"required,min=3,max=50"`
	DisplayName    string            `json:"display_name" binding:"required,min=1,max=100"`
	Password       string            `json:"password" binding:"required,min=8"`
	MFAEnabled     bool              `json:"mfa_enabled,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// UserUpdateRequest represents a request to update a user
type UserUpdateRequest struct {
	Email       string            `json:"email,omitempty" binding:"omitempty,email"`
	DisplayName string            `json:"display_name,omitempty" binding:"omitempty,min=1,max=100"`
	Password    string            `json:"password,omitempty" binding:"omitempty,min=8"`
	Status      UserStatus        `json:"status,omitempty" binding:"omitempty,oneof=active inactive pending suspended locked"`
	MFAEnabled  bool              `json:"mfa_enabled,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID               uuid.UUID         `json:"id"`
	OrganizationID   uuid.UUID         `json:"organization_id"`
	OrganizationName string            `json:"organization_name"`
	Email            string            `json:"email"`
	Username         string            `json:"username"`
	DisplayName      string            `json:"display_name"`
	Status           UserStatus        `json:"status"`
	MFAEnabled       bool              `json:"mfa_enabled"`
	LastLogin        *time.Time        `json:"last_login,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	RoleCount        int               `json:"role_count"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// UserFilter represents filters for listing users
type UserFilter struct {
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	Email          string    `json:"email,omitempty"`
	Username       string    `json:"username,omitempty"`
	Status         string    `json:"status,omitempty"`
	Page           int       `json:"page,omitempty"`
	Limit          int       `json:"limit,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	Email        string    `json:"email"`
	MFARequired  bool      `json:"mfa_required"`
}

// TokenRefreshRequest represents a token refresh request
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenRefreshResponse represents a token refresh response
type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
