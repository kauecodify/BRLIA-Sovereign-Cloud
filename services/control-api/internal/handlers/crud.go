package handlers

import (
	"net/http"
	"strconv"

	"github.com/brlia/control-api/internal/models"
	"github.com/brlia/control-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

// ListProjectsHandler lists projects
func ListProjectsHandler(svc *service.ProjectService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pagination(c)
		filter := &models.ProjectFilter{
			OrganizationID: parseUUIDQuery(c, "organization_id"),
			Name:           c.Query("name"),
			Status:         c.Query("status"),
			Page:           page,
			Limit:          limit,
		}

		projects, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects", "details": err.Error()})
			return
		}

		response := models.ProjectListResponse{Page: page, PageSize: limit, Total: total, TotalPages: int(total)/limit + 1}
		for _, project := range projects {
			resourceCount, _ := svc.CountResources(c.Request.Context(), project.ID)
			userCount, _ := svc.CountUsers(c.Request.Context(), project.ID)
			response.Projects = append(response.Projects, models.ProjectResponse{
				ID:             project.ID,
				OrganizationID: project.OrganizationID,
				Name:           project.Name,
				Description:    project.Description,
				Status:         project.Status,
				RegionID:       project.RegionID,
				CreatedAt:      project.CreatedAt,
				UpdatedAt:      project.UpdatedAt,
				Metadata:       project.Metadata,
				ResourceCount:  resourceCount,
				UserCount:      userCount,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateProjectHandler creates a new project
func CreateProjectHandler(svc *service.ProjectService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ProjectCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		project := &models.Project{
			OrganizationID: req.OrganizationID,
			Name:           req.Name,
			Description:    req.Description,
			RegionID:       req.RegionID,
			Metadata:       req.Metadata,
		}

		if err := svc.Create(c.Request.Context(), project); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.ProjectResponse{
			ID:             project.ID,
			OrganizationID: project.OrganizationID,
			Name:           project.Name,
			Description:    project.Description,
			Status:         project.Status,
			RegionID:       project.RegionID,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
			Metadata:       project.Metadata,
		})
	}
}

// GetProjectHandler retrieves a project by ID
func GetProjectHandler(svc *service.ProjectService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID", "details": err.Error()})
			return
		}

		project, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		c.JSON(http.StatusOK, models.ProjectResponse{
			ID:             project.ID,
			OrganizationID: project.OrganizationID,
			Name:           project.Name,
			Description:    project.Description,
			Status:         project.Status,
			RegionID:       project.RegionID,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
			Metadata:       project.Metadata,
		})
	}
}

// UpdateProjectHandler updates a project
func UpdateProjectHandler(svc *service.ProjectService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID", "details": err.Error()})
			return
		}

		project, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		var req models.ProjectUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		applyProjectUpdate(project, &req)

		if err := svc.Update(c.Request.Context(), project); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.ProjectResponse{
			ID:             project.ID,
			OrganizationID: project.OrganizationID,
			Name:           project.Name,
			Description:    project.Description,
			Status:         project.Status,
			RegionID:       project.RegionID,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
			Metadata:       project.Metadata,
		})
	}
}

// DeleteProjectHandler deletes a project
func DeleteProjectHandler(svc *service.ProjectService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID", "details": err.Error()})
			return
		}

		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ListUsersHandler lists users
func ListUsersHandler(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pagination(c)
		filter := &models.UserFilter{
			OrganizationID: parseUUIDQuery(c, "organization_id"),
			Email:          c.Query("email"),
			Username:       c.Query("username"),
			Status:         c.Query("status"),
			Page:           page,
			Limit:          limit,
		}

		users, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users", "details": err.Error()})
			return
		}

		response := models.UserListResponse{Page: page, PageSize: limit, Total: total, TotalPages: int(total)/limit + 1}
		for _, user := range users {
			roleCount, _ := svc.CountRoles(c.Request.Context(), user.ID)
			response.Users = append(response.Users, toUserResponse(user, roleCount))
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateUserHandler creates a new user
func CreateUserHandler(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UserCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		user := &models.User{
			OrganizationID: req.OrganizationID,
			Email:          req.Email,
			Username:       req.Username,
			DisplayName:    req.DisplayName,
			PasswordHash:   req.Password,
			MFAEnabled:     req.MFAEnabled,
			Metadata:       req.Metadata,
		}

		if err := svc.Create(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toUserResponse(user, 0))
	}
}

// GetUserHandler retrieves a user by ID
func GetUserHandler(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
			return
		}

		user, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		roleCount, _ := svc.CountRoles(c.Request.Context(), user.ID)
		c.JSON(http.StatusOK, toUserResponse(user, roleCount))
	}
}

// UpdateUserHandler updates a user
func UpdateUserHandler(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
			return
		}

		user, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var req models.UserUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		applyUserUpdate(user, &req)

		if err := svc.Update(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
			return
		}

		roleCount, _ := svc.CountRoles(c.Request.Context(), user.ID)
		c.JSON(http.StatusOK, toUserResponse(user, roleCount))
	}
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(svc *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
			return
		}

		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ListRolesHandler lists roles
func ListRolesHandler(svc *service.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pagination(c)
		filter := &models.RoleFilter{
			Name:     c.Query("name"),
			IsSystem: c.Query("is_system") == "true",
			Page:     page,
			Limit:    limit,
		}

		roles, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles", "details": err.Error()})
			return
		}

		response := models.RoleListResponse{Page: page, PageSize: limit, Total: total, TotalPages: int(total)/limit + 1}
		for _, role := range roles {
			response.Roles = append(response.Roles, models.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Permissions: role.Permissions,
				IsSystem:    role.IsSystem,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateRoleHandler creates a new role
func CreateRoleHandler(svc *service.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RoleCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		role := &models.Role{
			Name:        req.Name,
			Description: req.Description,
			Permissions: req.Permissions,
		}

		if err := svc.Create(c.Request.Context(), role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
}

// GetRoleHandler retrieves a role by ID
func GetRoleHandler(svc *service.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID", "details": err.Error()})
			return
		}

		role, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		c.JSON(http.StatusOK, models.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
}

// UpdateRoleHandler updates a role
func UpdateRoleHandler(svc *service.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID", "details": err.Error()})
			return
		}

		role, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		var req models.RoleUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		applyRoleUpdate(role, &req)

		if err := svc.Update(c.Request.Context(), role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
}

// DeleteRoleHandler deletes a role
func DeleteRoleHandler(svc *service.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID", "details": err.Error()})
			return
		}

		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ListQuotasHandler lists quotas
func ListQuotasHandler(svc *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pagination(c)
		filter := &models.QuotaFilter{
			ResourceType: c.Query("resource_type"),
			Scope:        c.Query("scope"),
			Page:         page,
			Limit:        limit,
		}

		quotas, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list quotas", "details": err.Error()})
			return
		}

		response := models.QuotaListResponse{Page: page, PageSize: limit, Total: total, TotalPages: int(total)/limit + 1}
		for _, quota := range quotas {
			response.Quotas = append(response.Quotas, toQuotaResponse(quota))
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateQuotaHandler creates a new quota
func CreateQuotaHandler(svc *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.QuotaCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		quota := &models.Quota{
			ResourceType: req.ResourceType,
			Unit:         req.Unit,
			Scope:        models.QuotaScope(req.Scope),
			Metadata:     req.Metadata,
		}

		if err := svc.Create(c.Request.Context(), quota); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quota", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toQuotaResponse(quota))
	}
}

// GetQuotaHandler retrieves a quota by ID
func GetQuotaHandler(svc *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quota ID", "details": err.Error()})
			return
		}

		quota, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quota not found"})
			return
		}

		c.JSON(http.StatusOK, toQuotaResponse(quota))
	}
}

// UpdateQuotaHandler updates a quota
func UpdateQuotaHandler(svc *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quota ID", "details": err.Error()})
			return
		}

		quota, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quota not found"})
			return
		}

		var req models.QuotaUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		applyQuotaUpdate(quota, &req)

		if err := svc.Update(c.Request.Context(), quota); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quota", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, toQuotaResponse(quota))
	}
}

// DeleteQuotaHandler deletes a quota
func DeleteQuotaHandler(svc *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quota ID", "details": err.Error()})
			return
		}

		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quota", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ListRegionsHandler lists regions
func ListRegionsHandler(svc *service.RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pagination(c)
		filter := &models.RegionFilter{
			Name:    c.Query("name"),
			Country: c.Query("country"),
			Status:  c.Query("status"),
			Page:    page,
			Limit:   limit,
		}

		regions, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list regions", "details": err.Error()})
			return
		}

		response := models.RegionListResponse{Page: page, PageSize: limit, Total: total, TotalPages: int(total)/limit + 1}
		for _, region := range regions {
			response.Regions = append(response.Regions, toRegionResponse(region))
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateRegionHandler creates a new region
func CreateRegionHandler(svc *service.RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegionCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		region := &models.Region{
			Name:        req.Name,
			Description: req.Description,
			Location:    req.Location,
			Country:     req.Country,
			IsDefault:   req.IsDefault,
			Metadata:    req.Metadata,
		}

		if err := svc.Create(c.Request.Context(), region); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create region", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, toRegionResponse(region))
	}
}

// GetRegionHandler retrieves a region by ID
func GetRegionHandler(svc *service.RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID", "details": err.Error()})
			return
		}

		region, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
			return
		}

		c.JSON(http.StatusOK, toRegionResponse(region))
	}
}

// UpdateRegionHandler updates a region
func UpdateRegionHandler(svc *service.RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID", "details": err.Error()})
			return
		}

		region, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
			return
		}

		var req models.RegionUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		applyRegionUpdate(region, &req)

		if err := svc.Update(c.Request.Context(), region); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update region", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, toRegionResponse(region))
	}
}

// DeleteRegionHandler deletes a region
func DeleteRegionHandler(svc *service.RegionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID", "details": err.Error()})
			return
		}

		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete region", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ListAuditLogsHandler lists audit logs
func ListAuditLogsHandler(svc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit := pagination(c)
		filter := &models.AuditLogFilter{
			Username:     c.Query("username"),
			Action:       c.Query("action"),
			ResourceType: c.Query("resource_type"),
			Status:       c.Query("status"),
			Page:         page,
			Limit:        limit,
		}

		logs, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list audit logs", "details": err.Error()})
			return
		}

		response := models.AuditLogListResponse{Page: page, PageSize: limit, Total: total, TotalPages: int(total)/limit + 1}
		for _, entry := range logs {
			response.AuditLogs = append(response.AuditLogs, toAuditLogResponse(entry))
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetAuditLogHandler retrieves an audit log by ID
func GetAuditLogHandler(svc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID", "details": err.Error()})
			return
		}

		entry, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
			return
		}

		c.JSON(http.StatusOK, toAuditLogResponse(entry))
	}
}

func parseUUIDQuery(c *gin.Context, key string) uuid.UUID {
	value := c.Query(key)
	if value == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func toUserResponse(user *models.User, roleCount int) models.UserResponse {
	return models.UserResponse{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Status:         user.Status,
		MFAEnabled:     user.MFAEnabled,
		LastLogin:      user.LastLogin,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Metadata:       user.Metadata,
		RoleCount:      roleCount,
	}
}

func toQuotaResponse(quota *models.Quota) models.QuotaResponse {
	var available int64
	if quota.Limit > quota.Used {
		available = quota.Limit - quota.Used
	}

	usagePercent := 0.0
	if quota.Limit > 0 {
		usagePercent = float64(quota.Used) / float64(quota.Limit) * 100
	}

	organizationName := ""
	if quota.OrganizationID != nil {
		organizationName = quota.OrganizationID.String()
	}

	projectName := ""
	if quota.ProjectID != nil {
		projectName = quota.ProjectID.String()
	}

	return models.QuotaResponse{
		ID:               quota.ID,
		OrganizationID:   quota.OrganizationID,
		OrganizationName: &organizationName,
		ProjectID:        quota.ProjectID,
		ProjectName:      &projectName,
		ResourceType:     quota.ResourceType,
		Limit:            quota.Limit,
		Used:             quota.Used,
		Available:        available,
		Unit:             quota.Unit,
		Scope:            quota.Scope,
		UsagePercent:     usagePercent,
		CreatedAt:        quota.CreatedAt,
		UpdatedAt:        quota.UpdatedAt,
		Metadata:         quota.Metadata,
	}
}

func toRegionResponse(region *models.Region) models.RegionResponse {
	return models.RegionResponse{
		ID:          region.ID,
		Name:        region.Name,
		Description: region.Description,
		Location:    region.Location,
		Country:     region.Country,
		Status:      region.Status,
		IsDefault:   region.IsDefault,
		CreatedAt:   region.CreatedAt,
		UpdatedAt:   region.UpdatedAt,
		Metadata:    region.Metadata,
	}
}

func toAuditLogResponse(entry *models.AuditLog) models.AuditLogResponse {
	return models.AuditLogResponse{
		ID:             entry.ID,
		OrganizationID: entry.OrganizationID,
		ProjectID:      entry.ProjectID,
		UserID:         entry.UserID,
		Username:       entry.Username,
		Action:         entry.Action,
		ResourceType:   entry.ResourceType,
		ResourceID:     entry.ResourceID,
		ResourceName:   entry.ResourceName,
		Details:        entry.Details,
		IPAddress:      entry.IPAddress,
		UserAgent:      entry.UserAgent,
		RequestID:      entry.RequestID,
		Status:         entry.Status,
		CreatedAt:      entry.CreatedAt,
		Metadata:       entry.Metadata,
	}
}

func applyProjectUpdate(project *models.Project, req *models.ProjectUpdateRequest) {
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = models.ProjectStatus(req.Status)
	}
	if req.Metadata != nil {
		project.Metadata = req.Metadata
	}
}

func applyUserUpdate(user *models.User, req *models.UserUpdateRequest) {
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Password != "" {
		user.PasswordHash = req.Password
	}
	if req.Status != "" {
		user.Status = models.UserStatus(req.Status)
	}
	user.MFAEnabled = req.MFAEnabled
	if req.Metadata != nil {
		user.Metadata = req.Metadata
	}
}

func applyRoleUpdate(role *models.Role, req *models.RoleUpdateRequest) {
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Permissions != nil {
		role.Permissions = req.Permissions
	}
}

func applyQuotaUpdate(quota *models.Quota, req *models.QuotaUpdateRequest) {
	quota.Limit = req.Limit
	if req.Metadata != nil {
		quota.Metadata = req.Metadata
	}
}

func applyRegionUpdate(region *models.Region, req *models.RegionUpdateRequest) {
	if req.Name != "" {
		region.Name = req.Name
	}
	if req.Description != "" {
		region.Description = req.Description
	}
	if req.Location != "" {
		region.Location = req.Location
	}
	if req.Country != "" {
		region.Country = req.Country
	}
	if req.Status != "" {
		region.Status = models.RegionStatus(req.Status)
	}
	region.IsDefault = req.IsDefault
	if req.Metadata != nil {
		region.Metadata = req.Metadata
	}
}
