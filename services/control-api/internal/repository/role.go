package repository

import (
	"context"
	"fmt"

	"github.com/brlia/control-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoleRepository handles database operations for roles
type RoleRepository struct {
	db *pgxpool.Pool
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	query := `
		INSERT INTO roles (id, name, description, permissions, is_system, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	args := []interface{}{
		role.ID,
		role.Name,
		role.Description,
		role.Permissions,
		role.IsSystem,
		role.CreatedAt,
		role.UpdatedAt,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := `
		SELECT id, name, description, permissions, is_system, created_at, updated_at, deleted_at
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL
	`

	var role models.Role
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Permissions,
		&role.IsSystem,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, permissions, is_system, created_at, updated_at, deleted_at
		FROM roles
		WHERE name = $1 AND deleted_at IS NULL
	`

	var role models.Role
	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Permissions,
		&role.IsSystem,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return &role, nil
}

// List retrieves a paginated list of roles
func (r *RoleRepository) List(ctx context.Context, filter *models.RoleFilter) ([]*models.Role, int64, error) {
	query := `
		SELECT id, name, description, permissions, is_system, created_at, updated_at, deleted_at
		FROM roles
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if filter.IsSystem {
		query += fmt.Sprintf(" AND is_system = $%d", argIndex)
		args = append(args, true)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:] + ")"
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.Permissions,
			&role.IsSystem,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, total, nil
}

// Update updates a role
func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	query := `
		UPDATE roles
		SET name = $2, description = $3, permissions = $4, is_system = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		role.ID,
		role.Name,
		role.Description,
		role.Permissions,
		role.IsSystem,
		role.UpdatedAt,
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

// Delete soft deletes a role
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE roles
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}
