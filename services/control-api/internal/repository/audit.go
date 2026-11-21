package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/brlia/control-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditRepository handles database operations for audit logs
type AuditRepository struct {
	db *pgxpool.Pool
}

// NewAuditRepository creates a new AuditRepository
func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(ctx context.Context, audit *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs 
		(id, organization_id, project_id, user_id, username, action, resource_type, resource_id, resource_name, details, ip_address, user_agent, request_id, status, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at
	`

	args := []interface{}{
		audit.ID,
		audit.OrganizationID,
		audit.ProjectID,
		audit.UserID,
		audit.Username,
		audit.Action,
		audit.ResourceType,
		audit.ResourceID,
		audit.ResourceName,
		audit.Details,
		audit.IPAddress,
		audit.UserAgent,
		audit.RequestID,
		audit.Status,
		audit.CreatedAt,
		audit.Metadata,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetByID retrieves an audit log by ID
func (r *AuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	query := `
		SELECT id, organization_id, project_id, user_id, username, action, resource_type, resource_id, resource_name, details, ip_address, user_agent, request_id, status, created_at, metadata
		FROM audit_logs
		WHERE id = $1
	`

	var audit models.AuditLog
	err := r.db.QueryRow(ctx, query, id).Scan(
		&audit.ID,
		&audit.OrganizationID,
		&audit.ProjectID,
		&audit.UserID,
		&audit.Username,
		&audit.Action,
		&audit.ResourceType,
		&audit.ResourceID,
		&audit.ResourceName,
		&audit.Details,
		&audit.IPAddress,
		&audit.UserAgent,
		&audit.RequestID,
		&audit.Status,
		&audit.CreatedAt,
		&audit.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}

	return &audit, nil
}

// List retrieves a paginated list of audit logs
func (r *AuditRepository) List(ctx context.Context, filter *models.AuditLogFilter) ([]*models.AuditLog, int64, error) {
	query := `
		SELECT id, organization_id, project_id, user_id, username, action, resource_type, resource_id, resource_name, details, ip_address, user_agent, request_id, status, created_at, metadata
		FROM audit_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filter.OrganizationID != nil {
		query += fmt.Sprintf(" AND organization_id = $%d", argIndex)
		args = append(args, *filter.OrganizationID)
		argIndex++
	}

	if filter.ProjectID != nil {
		query += fmt.Sprintf(" AND project_id = $%d", argIndex)
		args = append(args, *filter.ProjectID)
		argIndex++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Username != "" {
		query += fmt.Sprintf(" AND username = $%d", argIndex)
		args = append(args, filter.Username)
		argIndex++
	}

	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, filter.Action)
		argIndex++
	}

	if filter.ResourceType != "" {
		query += fmt.Sprintf(" AND resource_type = $%d", argIndex)
		args = append(args, filter.ResourceType)
		argIndex++
	}

	if filter.ResourceID != nil {
		query += fmt.Sprintf(" AND resource_id = $%d", argIndex)
		args = append(args, *filter.ResourceID)
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.StartTime != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.StartTime)
		argIndex++
	}

	if filter.EndTime != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.EndTime)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM audit_logs WHERE 1=1"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:] + ")"
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Add pagination and ordering
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.Limit, offset)
	} else {
		query += " ORDER BY created_at DESC"
	}

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var audits []*models.AuditLog
	for rows.Next() {
		var audit models.AuditLog
		err := rows.Scan(
			&audit.ID,
			&audit.OrganizationID,
			&audit.ProjectID,
			&audit.UserID,
			&audit.Username,
			&audit.Action,
			&audit.ResourceType,
			&audit.ResourceID,
			&audit.ResourceName,
			&audit.Details,
			&audit.IPAddress,
			&audit.UserAgent,
			&audit.RequestID,
			&audit.Status,
			&audit.CreatedAt,
			&audit.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}
		audits = append(audits, &audit)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating audit logs: %w", err)
	}

	return audits, total, nil
}

// GetSummary returns a summary of audit logs
func (r *AuditRepository) GetSummary(ctx context.Context, startTime, endTime *time.Time) (*models.AuditSummary, error) {
	// This would implement summary queries
	// For now, return a simple summary
	return &models.AuditSummary{
		TotalActions:    0,
		ActionsByType:   make(map[string]int64),
		ActionsByUser:   make(map[string]int64),
		ActionsByStatus: make(map[string]int64),
		RecentActions:   []models.AuditLogResponse{},
	}, nil
}
