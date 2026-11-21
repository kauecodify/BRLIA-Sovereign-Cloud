package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/brlia/control-api/internal/models"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new ProjectRepository
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (id, organization_id, name, description, status, region_id, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	args := []interface{}{
		project.ID,
		project.OrganizationID,
		project.Name,
		project.Description,
		project.Status,
		project.RegionID,
		project.Metadata,
		project.CreatedAt,
		project.UpdatedAt,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `
		SELECT id, organization_id, name, description, status, region_id, created_at, updated_at, deleted_at, metadata
		FROM projects
		WHERE id = $1 AND deleted_at IS NULL
	`

	var project models.Project
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Description,
		&project.Status,
		&project.RegionID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.DeletedAt,
		&project.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// List retrieves a paginated list of projects
func (r *ProjectRepository) List(ctx context.Context, filter *models.ProjectFilter) ([]*models.Project, int64, error) {
	query := `
		SELECT id, organization_id, name, description, status, region_id, created_at, updated_at, deleted_at, metadata
		FROM projects
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if filter.OrganizationID != uuid.Nil {
		query += fmt.Sprintf(" AND organization_id = $%d", argIndex)
		args = append(args, filter.OrganizationID)
		argIndex++
	}

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

	if filter.RegionID != uuid.Nil {
		query += fmt.Sprintf(" AND region_id = $%d", argIndex)
		args = append(args, filter.RegionID)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:] + ")"
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.OrganizationID,
			&project.Name,
			&project.Description,
			&project.Status,
			&project.RegionID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.DeletedAt,
			&project.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, total, nil
}

// Update updates a project
func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET name = $2, description = $3, status = $4, region_id = $5, metadata = $6, updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		project.ID,
		project.Name,
		project.Description,
		project.Status,
		project.RegionID,
		project.Metadata,
		project.UpdatedAt,
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

// Delete soft deletes a project
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE projects
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// CountResources counts the number of resources in a project
func (r *ProjectRepository) CountResources(ctx context.Context, projectID uuid.UUID) (int, error) {
	// This would count various resource types
	// For now, return 0 as a placeholder
	return 0, nil
}

// CountUsers counts the number of users with access to a project
func (r *ProjectRepository) CountUsers(ctx context.Context, projectID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN project_users pu ON u.id = pu.user_id
		WHERE pu.project_id = $1 AND u.deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(ctx, query, projectID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count project users: %w", err)
	}

	return count, nil
}
