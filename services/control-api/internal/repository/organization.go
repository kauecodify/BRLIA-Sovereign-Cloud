package repository

import (
	"context"
	"fmt"

	"github.com/brlia/control-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrganizationRepository handles database operations for organizations
type OrganizationRepository struct {
	db *pgxpool.Pool
}

// NewOrganizationRepository creates a new OrganizationRepository
func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *models.Organization) error {
	query := `
		INSERT INTO organizations (id, name, description, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	args := []interface{}{
		org.ID,
		org.Name,
		org.Description,
		org.Status,
		org.Metadata,
		org.CreatedAt,
		org.UpdatedAt,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

// GetByID retrieves an organization by ID
func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at, metadata
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`

	var org models.Organization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Description,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
		&org.DeletedAt,
		&org.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return &org, nil
}

// GetByName retrieves an organization by name
func (r *OrganizationRepository) GetByName(ctx context.Context, name string) (*models.Organization, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at, metadata
		FROM organizations
		WHERE name = $1 AND deleted_at IS NULL
	`

	var org models.Organization
	err := r.db.QueryRow(ctx, query, name).Scan(
		&org.ID,
		&org.Name,
		&org.Description,
		&org.Status,
		&org.CreatedAt,
		&org.UpdatedAt,
		&org.DeletedAt,
		&org.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get organization by name: %w", err)
	}

	return &org, nil
}

// List retrieves a paginated list of organizations
func (r *OrganizationRepository) List(ctx context.Context, filter *models.OrganizationFilter) ([]*models.Organization, int64, error) {
	// Build query
	query := `
		SELECT id, name, description, status, created_at, updated_at, deleted_at, metadata
		FROM organizations
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM organizations WHERE deleted_at IS NULL"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:]
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var organizations []*models.Organization
	for rows.Next() {
		var org models.Organization
		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.Description,
			&org.Status,
			&org.CreatedAt,
			&org.UpdatedAt,
			&org.DeletedAt,
			&org.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization: %w", err)
		}
		organizations = append(organizations, &org)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating organizations: %w", err)
	}

	return organizations, total, nil
}

// Update updates an organization
func (r *OrganizationRepository) Update(ctx context.Context, org *models.Organization) error {
	query := `
		UPDATE organizations
		SET name = $2, description = $3, status = $4, metadata = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		org.ID,
		org.Name,
		org.Description,
		org.Status,
		org.Metadata,
		org.UpdatedAt,
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	return nil
}

// Delete soft deletes an organization
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE organizations
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

// HardDelete permanently deletes an organization
func (r *OrganizationRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete organization: %w", err)
	}

	return nil
}

// CountProjects counts the number of projects in an organization
func (r *OrganizationRepository) CountProjects(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM projects
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count projects: %w", err)
	}

	return count, nil
}

// CountUsers counts the number of users in an organization
func (r *OrganizationRepository) CountUsers(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
