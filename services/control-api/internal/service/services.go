package service

import (
	"context"
	"log"
	"time"

	"github.com/brlia/control-api/internal/models"
	"github.com/brlia/control-api/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func writeAudit(ctx context.Context, auditRepo *repository.AuditRepository, action models.AuditAction, resourceType string, resourceID uuid.UUID, resourceName string, details map[string]string) {
	if auditRepo == nil {
		return
	}

	audit := &models.AuditLog{
		ID:           uuid.New(),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		Status:       models.AuditStatusCompleted,
		CreatedAt:    time.Now(),
	}

	if err := auditRepo.Create(ctx, audit); err != nil {
		log.Printf("Failed to persist audit log for %s %s: %v", resourceType, resourceID, err)
	}
}

// ProjectService provides business logic for projects
type ProjectService struct {
	repo      *repository.ProjectRepository
	orgRepo   *repository.OrganizationRepository
	auditRepo *repository.AuditRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService(repo *repository.ProjectRepository, orgRepo *repository.OrganizationRepository, auditRepo *repository.AuditRepository) *ProjectService {
	return &ProjectService{repo: repo, orgRepo: orgRepo, auditRepo: auditRepo}
}

// Create creates a new project
func (s *ProjectService) Create(ctx context.Context, project *models.Project) error {
	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	if project.UpdatedAt.IsZero() {
		project.UpdatedAt = time.Now()
	}
	if project.Status == "" {
		project.Status = models.ProjectStatusPending
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionCreate, "project", project.ID, project.Name, nil)
	return nil
}

// GetByID retrieves a project by ID
func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves a paginated list of projects
func (s *ProjectService) List(ctx context.Context, filter *models.ProjectFilter) ([]*models.Project, int64, error) {
	return s.repo.List(ctx, filter)
}

// Update updates a project
func (s *ProjectService) Update(ctx context.Context, project *models.Project) error {
	project.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, project); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionUpdate, "project", project.ID, project.Name, nil)
	return nil
}

// Delete deletes a project
func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// CountResources counts resources in a project
func (s *ProjectService) CountResources(ctx context.Context, projectID uuid.UUID) (int, error) {
	return s.repo.CountResources(ctx, projectID)
}

// CountUsers counts users in a project
func (s *ProjectService) CountUsers(ctx context.Context, projectID uuid.UUID) (int, error) {
	return s.repo.CountUsers(ctx, projectID)
}

// UserService provides business logic for users
type UserService struct {
	repo      *repository.UserRepository
	auditRepo *repository.AuditRepository
}

// NewUserService creates a new UserService
func NewUserService(repo *repository.UserRepository, auditRepo *repository.AuditRepository) *UserService {
	return &UserService{repo: repo, auditRepo: auditRepo}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}
	if user.Status == "" {
		user.Status = models.UserStatusPending
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionCreate, "user", user.ID, user.Username, nil)
	return nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// List retrieves a paginated list of users
func (s *UserService) List(ctx context.Context, filter *models.UserFilter) ([]*models.User, int64, error) {
	return s.repo.List(ctx, filter)
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	if user.PasswordHash != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hash)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionUpdate, "user", user.ID, user.Username, nil)
	return nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// CountRoles counts roles assigned to a user
func (s *UserService) CountRoles(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.CountRoles(ctx, userID)
}

// RoleService provides business logic for roles
type RoleService struct {
	repo      *repository.RoleRepository
	auditRepo *repository.AuditRepository
}

// NewRoleService creates a new RoleService
func NewRoleService(repo *repository.RoleRepository, auditRepo *repository.AuditRepository) *RoleService {
	return &RoleService{repo: repo, auditRepo: auditRepo}
}

// Create creates a new role
func (s *RoleService) Create(ctx context.Context, role *models.Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	if role.CreatedAt.IsZero() {
		role.CreatedAt = time.Now()
	}
	if role.UpdatedAt.IsZero() {
		role.UpdatedAt = time.Now()
	}

	if err := s.repo.Create(ctx, role); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionCreate, "role", role.ID, role.Name, nil)
	return nil
}

// GetByID retrieves a role by ID
func (s *RoleService) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByName retrieves a role by name
func (s *RoleService) GetByName(ctx context.Context, name string) (*models.Role, error) {
	return s.repo.GetByName(ctx, name)
}

// List retrieves a paginated list of roles
func (s *RoleService) List(ctx context.Context, filter *models.RoleFilter) ([]*models.Role, int64, error) {
	return s.repo.List(ctx, filter)
}

// Update updates a role
func (s *RoleService) Update(ctx context.Context, role *models.Role) error {
	role.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, role); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionUpdate, "role", role.ID, role.Name, nil)
	return nil
}

// Delete deletes a role
func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// QuotaService provides business logic for quotas
type QuotaService struct {
	repo      *repository.QuotaRepository
	orgRepo   *repository.OrganizationRepository
	projRepo  *repository.ProjectRepository
	auditRepo *repository.AuditRepository
}

// NewQuotaService creates a new QuotaService
func NewQuotaService(repo *repository.QuotaRepository, orgRepo *repository.OrganizationRepository, projRepo *repository.ProjectRepository, auditRepo *repository.AuditRepository) *QuotaService {
	return &QuotaService{repo: repo, orgRepo: orgRepo, projRepo: projRepo, auditRepo: auditRepo}
}

// Create creates a new quota
func (s *QuotaService) Create(ctx context.Context, quota *models.Quota) error {
	if quota.ID == uuid.Nil {
		quota.ID = uuid.New()
	}
	if quota.CreatedAt.IsZero() {
		quota.CreatedAt = time.Now()
	}
	if quota.UpdatedAt.IsZero() {
		quota.UpdatedAt = time.Now()
	}

	if err := s.repo.Create(ctx, quota); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionCreate, "quota", quota.ID, quota.ResourceType, nil)
	return nil
}

// GetByID retrieves a quota by ID
func (s *QuotaService) GetByID(ctx context.Context, id uuid.UUID) (*models.Quota, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves a paginated list of quotas
func (s *QuotaService) List(ctx context.Context, filter *models.QuotaFilter) ([]*models.Quota, int64, error) {
	return s.repo.List(ctx, filter)
}

// Update updates a quota
func (s *QuotaService) Update(ctx context.Context, quota *models.Quota) error {
	quota.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, quota); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionUpdate, "quota", quota.ID, quota.ResourceType, nil)
	return nil
}

// Delete deletes a quota
func (s *QuotaService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// RegionService provides business logic for regions
type RegionService struct {
	repo      *repository.RegionRepository
	auditRepo *repository.AuditRepository
}

// NewRegionService creates a new RegionService
func NewRegionService(repo *repository.RegionRepository, auditRepo *repository.AuditRepository) *RegionService {
	return &RegionService{repo: repo, auditRepo: auditRepo}
}

// Create creates a new region
func (s *RegionService) Create(ctx context.Context, region *models.Region) error {
	if region.ID == uuid.Nil {
		region.ID = uuid.New()
	}
	if region.CreatedAt.IsZero() {
		region.CreatedAt = time.Now()
	}
	if region.UpdatedAt.IsZero() {
		region.UpdatedAt = time.Now()
	}
	if region.Status == "" {
		region.Status = models.RegionStatusActive
	}

	if err := s.repo.Create(ctx, region); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionCreate, "region", region.ID, region.Name, nil)
	return nil
}

// GetByID retrieves a region by ID
func (s *RegionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Region, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByName retrieves a region by name
func (s *RegionService) GetByName(ctx context.Context, name string) (*models.Region, error) {
	return s.repo.GetByName(ctx, name)
}

// List retrieves a paginated list of regions
func (s *RegionService) List(ctx context.Context, filter *models.RegionFilter) ([]*models.Region, int64, error) {
	return s.repo.List(ctx, filter)
}

// Update updates a region
func (s *RegionService) Update(ctx context.Context, region *models.Region) error {
	region.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, region); err != nil {
		return err
	}

	writeAudit(ctx, s.auditRepo, models.AuditActionUpdate, "region", region.ID, region.Name, nil)
	return nil
}

// Delete deletes a region
func (s *RegionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// AuditService provides read access to audit logs
type AuditService struct {
	repo *repository.AuditRepository
}

// NewAuditService creates a new AuditService
func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// List retrieves a paginated list of audit logs
func (s *AuditService) List(ctx context.Context, filter *models.AuditLogFilter) ([]*models.AuditLog, int64, error) {
	return s.repo.List(ctx, filter)
}

// GetByID retrieves an audit log by ID
func (s *AuditService) GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	return s.repo.GetByID(ctx, id)
}

// GetSummary aggregates audit logs for a period
func (s *AuditService) GetSummary(ctx context.Context, startTime, endTime *time.Time) (*models.AuditSummary, error) {
	return s.repo.GetSummary(ctx, startTime, endTime)
}
