package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/brlia/control-api/internal/models"
)

// RegionRepository handles database operations for regions
type RegionRepository struct {
	db *pgxpool.Pool
}

// NewRegionRepository creates a new RegionRepository
func NewRegionRepository(db *pgxpool.Pool) *RegionRepository {
	return &RegionRepository{db: db}
}

// Create creates a new region
func (r *RegionRepository) Create(ctx context.Context, region *models.Region) error {
	query := `
		INSERT INTO regions (id, name, description, location, country, status, is_default, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	args := []interface{}{
		region.ID,
		region.Name,
		region.Description,
		region.Location,
		region.Country,
		region.Status,
		region.IsDefault,
		region.Metadata,
		region.CreatedAt,
		region.UpdatedAt,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create region: %w", err)
	}

	return nil
}

// GetByID retrieves a region by ID
func (r *RegionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Region, error) {
	query := `
		SELECT id, name, description, location, country, status, is_default, created_at, updated_at, deleted_at, metadata
		FROM regions
		WHERE id = $1 AND deleted_at IS NULL
	`

	var region models.Region
	err := r.db.QueryRow(ctx, query, id).Scan(
		&region.ID,
		&region.Name,
		&region.Description,
		&region.Location,
		&region.Country,
		&region.Status,
		&region.IsDefault,
		&region.CreatedAt,
		&region.UpdatedAt,
		&region.DeletedAt,
		&region.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get region: %w", err)
	}

	return &region, nil
}

// GetByName retrieves a region by name
func (r *RegionRepository) GetByName(ctx context.Context, name string) (*models.Region, error) {
	query := `
		SELECT id, name, description, location, country, status, is_default, created_at, updated_at, deleted_at, metadata
		FROM regions
		WHERE name = $1 AND deleted_at IS NULL
	`

	var region models.Region
	err := r.db.QueryRow(ctx, query, name).Scan(
		&region.ID,
		&region.Name,
		&region.Description,
		&region.Location,
		&region.Country,
		&region.Status,
		&region.IsDefault,
		&region.CreatedAt,
		&region.UpdatedAt,
		&region.DeletedAt,
		&region.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get region by name: %w", err)
	}

	return &region, nil
}

// List retrieves a paginated list of regions
func (r *RegionRepository) List(ctx context.Context, filter *models.RegionFilter) ([]*models.Region, int64, error) {
	query := `
		SELECT id, name, description, location, country, status, is_default, created_at, updated_at, deleted_at, metadata
		FROM regions
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country = $%d", argIndex)
		args = append(args, filter.Country)
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM regions WHERE deleted_at IS NULL"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:] + ")"
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count regions: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list regions: %w", err)
	}
	defer rows.Close()

	var regions []*models.Region
	for rows.Next() {
		var region models.Region
		err := rows.Scan(
			&region.ID,
			&region.Name,
			&region.Description,
			&region.Location,
			&region.Country,
			&region.Status,
			&region.IsDefault,
			&region.CreatedAt,
			&region.UpdatedAt,
			&region.DeletedAt,
			&region.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan region: %w", err)
		}
		regions = append(regions, &region)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating regions: %w", err)
	}

	return regions, total, nil
}

// Update updates a region
func (r *RegionRepository) Update(ctx context.Context, region *models.Region) error {
	query := `
		UPDATE regions
		SET name = $2, description = $3, location = $4, country = $5, status = $6, is_default = $7, metadata = $8, updated_at = $9
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		region.ID,
		region.Name,
		region.Description,
		region.Location,
		region.Country,
		region.Status,
		region.IsDefault,
		region.Metadata,
		region.UpdatedAt,
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update region: %w", err)
	}

	return nil
}

// Delete soft deletes a region
func (r *RegionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE regions
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete region: %w", err)
	}

	return nil
}
