package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/brlia/control-api/internal/models"
)

// QuotaRepository handles database operations for quotas
type QuotaRepository struct {
	db *pgxpool.Pool
}

// NewQuotaRepository creates a new QuotaRepository
func NewQuotaRepository(db *pgxpool.Pool) *QuotaRepository {
	return &QuotaRepository{db: db}
}

// Create creates a new quota
func (r *QuotaRepository) Create(ctx context.Context, quota *models.Quota) error {
	query := `
		INSERT INTO quotas (id, organization_id, project_id, resource_type, limit, used, unit, scope, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at, updated_at
	`

	args := []interface{}{
		quota.ID,
		quota.OrganizationID,
		quota.ProjectID,
		quota.ResourceType,
		quota.Limit,
		quota.Used,
		quota.Unit,
		quota.Scope,
		quota.Metadata,
		quota.CreatedAt,
		quota.UpdatedAt,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create quota: %w", err)
	}

	return nil
}

// GetByID retrieves a quota by ID
func (r *QuotaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Quota, error) {
	query := `
		SELECT id, organization_id, project_id, resource_type, limit, used, unit, scope, created_at, updated_at, deleted_at, metadata
		FROM quotas
		WHERE id = $1 AND deleted_at IS NULL
	`

	var quota models.Quota
	err := r.db.QueryRow(ctx, query, id).Scan(
		&quota.ID,
		&quota.OrganizationID,
		&quota.ProjectID,
		&quota.ResourceType,
		&quota.Limit,
		&quota.Used,
		&quota.Unit,
		&quota.Scope,
		&quota.CreatedAt,
		&quota.UpdatedAt,
		&quota.DeletedAt,
		&quota.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get quota: %w", err)
	}

	return &quota, nil
}

// List retrieves a paginated list of quotas
func (r *QuotaRepository) List(ctx context.Context, filter *models.QuotaFilter) ([]*models.Quota, int64, error) {
	query := `
		SELECT id, organization_id, project_id, resource_type, limit, used, unit, scope, created_at, updated_at, deleted_at, metadata
		FROM quotas
		WHERE deleted_at IS NULL
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

	if filter.ResourceType != "" {
		query += fmt.Sprintf(" AND resource_type = $%d", argIndex)
		args = append(args, filter.ResourceType)
		argIndex++
	}

	if filter.Scope != "" {
		query += fmt.Sprintf(" AND scope = $%d", argIndex)
		args = append(args, filter.Scope)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM quotas WHERE deleted_at IS NULL"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:] + ")"
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count quotas: %w", err)
	}

	// Add pagination
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
		return nil, 0, fmt.Errorf("failed to list quotas: %w", err)
	}
	defer rows.Close()

	var quotas []*models.Quota
	for rows.Next() {
		var quota models.Quota
		err := rows.Scan(
			&quota.ID,
			&quota.OrganizationID,
			&quota.ProjectID,
			&quota.ResourceType,
			&quota.Limit,
			&quota.Used,
			&quota.Unit,
			&quota.Scope,
			&quota.CreatedAt,
			&quota.UpdatedAt,
			&quota.DeletedAt,
			&quota.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quota: %w", err)
		}
		quotas = append(quotas, &quota)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating quotas: %w", err)
	}

	return quotas, total, nil
}

// Update updates a quota
func (r *QuotaRepository) Update(ctx context.Context, quota *models.Quota) error {
	query := `
		UPDATE quotas
		SET limit = $2, used = $3, metadata = $4, updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		quota.ID,
		quota.Limit,
		quota.Used,
		quota.Metadata,
		quota.UpdatedAt,
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update quota: %w", err)
	}

	return nil
}

// Delete soft deletes a quota
func (r *QuotaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE quotas
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quota: %w", err)
	}

	return nil
}
