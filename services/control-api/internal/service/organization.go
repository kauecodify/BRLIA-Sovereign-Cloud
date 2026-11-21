package service

import (
	"context"
	"log"
	"time"

	"github.com/brlia/control-api/internal/models"
	"github.com/brlia/control-api/internal/repository"
	"github.com/google/uuid"
)

// OrganizationService provides business logic for organizations
type OrganizationService struct {
	repo      *repository.OrganizationRepository
	auditRepo *repository.AuditRepository
}

// NewOrganizationService creates a new OrganizationService
func NewOrganizationService(repo *repository.OrganizationRepository, auditRepo *repository.AuditRepository) *OrganizationService {
	return &OrganizationService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

// Create creates a new organization
func (s *OrganizationService) Create(ctx context.Context, org *models.Organization) error {
	// Set timestamps
	if org.CreatedAt.IsZero() {
		org.CreatedAt = time.Now()
	}
	if org.UpdatedAt.IsZero() {
		org.UpdatedAt = time.Now()
	}

	// Set default status
	if org.Status == "" {
		org.Status = models.OrganizationStatusPending
	}

	// Create in database
	if err := s.repo.Create(ctx, org); err != nil {
		return err
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "organization", org.ID, org.Name, nil)

	return nil
}

// GetByID retrieves an organization by ID
func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionRead, "organization", org.ID, org.Name, nil)

	return org, nil
}

// List retrieves a paginated list of organizations
func (s *OrganizationService) List(ctx context.Context, filter *models.OrganizationFilter) ([]*models.Organization, int64, error) {
	orgs, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionRead, "organization", uuid.Nil, "list", nil)

	return orgs, total, nil
}

// Update updates an organization
func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) error {
	// Set updated timestamp
	org.UpdatedAt = time.Now()

	// Update in database
	if err := s.repo.Update(ctx, org); err != nil {
		return err
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "organization", org.ID, org.Name, nil)

	return nil
}

// Delete deletes an organization
func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	// Get organization for audit
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "organization", org.ID, org.Name, nil)

	return nil
}

// CountProjects counts the number of projects in an organization
func (s *OrganizationService) CountProjects(ctx context.Context, orgID uuid.UUID) (int, error) {
	return s.repo.CountProjects(ctx, orgID)
}

// CountUsers counts the number of users in an organization
func (s *OrganizationService) CountUsers(ctx context.Context, orgID uuid.UUID) (int, error) {
	return s.repo.CountUsers(ctx, orgID)
}

// logAudit logs an audit entry
func (s *OrganizationService) logAudit(ctx context.Context, action models.AuditAction, resourceType string, resourceID uuid.UUID, resourceName string, details map[string]string) {
	if s.auditRepo == nil {
		return
	}

	audit := &models.AuditLog{
		ID:           uuid.New(),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		Details:      "",
		Status:       models.AuditStatusCompleted,
		CreatedAt:    time.Now(),
		Metadata:     details,
	}

	if err := s.auditRepo.Create(ctx, audit); err != nil {
		log.Printf("Failed to persist audit log for organization %s: %v", resourceID, err)
	}
}
