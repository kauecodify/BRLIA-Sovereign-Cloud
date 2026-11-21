package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/brlia/control-api/internal/models"
	"github.com/brlia/control-api/internal/service"
)

// ListOrganizationsHandler lists all organizations
func ListOrganizationsHandler(svc *service.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		name := c.Query("name")
		status := c.Query("status")

		filter := &models.OrganizationFilter{
			Page:   page,
			Limit:  limit,
			Name:   name,
			Status: status,
		}

		// Get organizations
		organizations, total, err := svc.List(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to list organizations",
				"details": err.Error(),
			})
			return
		}

		// Build response
		response := models.OrganizationListResponse{
			Page:       page,
			PageSize:   limit,
			Total:      total,
			TotalPages: int(total)/limit + 1,
		}

		for _, org := range organizations {
			// Get counts
			projectCount, _ := svc.CountProjects(c.Request.Context(), org.ID)
			userCount, _ := svc.CountUsers(c.Request.Context(), org.ID)

			response.Organizations = append(response.Organizations, models.OrganizationResponse{
				ID:          org.ID,
				Name:        org.Name,
				Description: org.Description,
				Status:      org.Status,
				CreatedAt:   org.CreatedAt,
				UpdatedAt:   org.UpdatedAt,
				Metadata:    org.Metadata,
				ProjectCount: projectCount,
				UserCount:    userCount,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// CreateOrganizationHandler creates a new organization
func CreateOrganizationHandler(svc *service.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.OrganizationCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Create organization
		org := &models.Organization{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: req.Description,
			Status:      models.OrganizationStatusPending,
			Metadata:    req.Metadata,
		}

		if err := svc.Create(c.Request.Context(), org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create organization",
				"details": err.Error(),
			})
			return
		}

		// Build response
		response := models.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			Status:      org.Status,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
			Metadata:    org.Metadata,
		}

		c.JSON(http.StatusCreated, response)
	}
}

// GetOrganizationHandler retrieves an organization by ID
func GetOrganizationHandler(svc *service.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid organization ID",
				"details": err.Error(),
			})
			return
		}

		// Get organization
		org, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}

		// Get counts
		projectCount, _ := svc.CountProjects(c.Request.Context(), org.ID)
		userCount, _ := svc.CountUsers(c.Request.Context(), org.ID)

		// Build response
		response := models.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			Status:      org.Status,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
			Metadata:    org.Metadata,
			ProjectCount: projectCount,
			UserCount:    userCount,
		}

		c.JSON(http.StatusOK, response)
	}
}

// UpdateOrganizationHandler updates an organization
func UpdateOrganizationHandler(svc *service.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid organization ID",
				"details": err.Error(),
			})
			return
		}

		var req models.OrganizationUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Get existing organization
		org, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Organization not found",
			})
			return
		}

		// Update fields
		if req.Name != "" {
			org.Name = req.Name
		}
		if req.Description != "" {
			org.Description = req.Description
		}
		if req.Status != "" {
			org.Status = req.Status
		}
		if req.Metadata != nil {
			org.Metadata = req.Metadata
		}

		if err := svc.Update(c.Request.Context(), org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update organization",
				"details": err.Error(),
			})
			return
		}

		// Get counts
		projectCount, _ := svc.CountProjects(c.Request.Context(), org.ID)
		userCount, _ := svc.CountUsers(c.Request.Context(), org.ID)

		// Build response
		response := models.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			Status:      org.Status,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
			Metadata:    org.Metadata,
			ProjectCount: projectCount,
			UserCount:    userCount,
		}

		c.JSON(http.StatusOK, response)
	}
}

// DeleteOrganizationHandler deletes an organization
func DeleteOrganizationHandler(svc *service.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid organization ID",
				"details": err.Error(),
			})
			return
		}

		// Delete organization
		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete organization",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
