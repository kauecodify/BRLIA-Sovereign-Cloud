package repository

import (
	"context"
	"fmt"

	"github.com/brlia/control-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, organization_id, email, username, display_name, password_hash, status, mfa_enabled, mfa_secret, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at
	`

	args := []interface{}{
		user.ID,
		user.OrganizationID,
		user.Email,
		user.Username,
		user.DisplayName,
		user.PasswordHash,
		user.Status,
		user.MFAEnabled,
		user.MFASecret,
		user.CreatedAt,
		user.UpdatedAt,
		user.Metadata,
	}

	_, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, organization_id, email, username, display_name, password_hash, status, mfa_enabled, mfa_secret, last_login, created_at, updated_at, deleted_at, metadata
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.Username,
		&user.DisplayName,
		&user.PasswordHash,
		&user.Status,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, organization_id, email, username, display_name, password_hash, status, mfa_enabled, mfa_secret, last_login, created_at, updated_at, deleted_at, metadata
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.Username,
		&user.DisplayName,
		&user.PasswordHash,
		&user.Status,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, organization_id, email, username, display_name, password_hash, status, mfa_enabled, mfa_secret, last_login, created_at, updated_at, deleted_at, metadata
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.Username,
		&user.DisplayName,
		&user.PasswordHash,
		&user.Status,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Metadata,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// List retrieves a paginated list of users
func (r *UserRepository) List(ctx context.Context, filter *models.UserFilter) ([]*models.User, int64, error) {
	query := `
		SELECT id, organization_id, email, username, display_name, password_hash, status, mfa_enabled, mfa_secret, last_login, created_at, updated_at, deleted_at, metadata
		FROM users
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if filter.OrganizationID != uuid.Nil {
		query += fmt.Sprintf(" AND organization_id = $%d", argIndex)
		args = append(args, filter.OrganizationID)
		argIndex++
	}

	if filter.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Email+"%")
		argIndex++
	}

	if filter.Username != "" {
		query += fmt.Sprintf(" AND username ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Username+"%")
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	if len(args) > 0 {
		countQuery += " AND (" + query[len("SELECT ")+7:] + ")"
	}

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.OrganizationID,
			&user.Email,
			&user.Username,
			&user.DisplayName,
			&user.PasswordHash,
			&user.Status,
			&user.MFAEnabled,
			&user.MFASecret,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $2, display_name = $3, password_hash = $4, status = $5, mfa_enabled = $6, mfa_secret = $7, last_login = $8, updated_at = $9, metadata = $10
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		user.ID,
		user.Email,
		user.DisplayName,
		user.PasswordHash,
		user.Status,
		user.MFAEnabled,
		user.MFASecret,
		user.LastLogin,
		user.UpdatedAt,
		user.Metadata,
	}

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login time for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET last_login = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// CountRoles counts the number of roles assigned to a user
func (r *UserRepository) CountRoles(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_roles
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user roles: %w", err)
	}

	return count, nil
}
