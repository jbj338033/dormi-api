package handler

import (
	"encoding/json"
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// GetAll godoc
// @Summary 감사 로그 조회
// @Description 감사 로그 목록 조회 (관리자 전용)
// @Tags 감사로그
// @Produce json
// @Security BearerAuth
// @Param userId query string false "사용자 ID"
// @Param action query string false "작업 유형"
// @Param entityType query string false "엔티티 유형"
// @Param startDate query string false "시작일 (RFC3339)"
// @Param endDate query string false "종료일 (RFC3339)"
// @Param page query int false "페이지" default(1)
// @Param limit query int false "페이지당 개수" default(20)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.AuditLogResponse}
// @Failure 400 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Router /audit-logs [get]
func (h *AuditHandler) GetAll(c *gin.Context) {
	var query dto.AuditQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	logs, total, err := h.auditService.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.AuditLogResponse{}
	for _, log := range logs {
		responses = append(responses, toAuditLogResponse(&log))
	}

	totalPages := int(total) / query.Limit
	if int(total)%query.Limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Data:    responses,
		Meta: &dto.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func toAuditLogResponse(log *model.AuditLog) dto.AuditLogResponse {
	resp := dto.AuditLogResponse{
		ID:         log.ID,
		Action:     string(log.Action),
		EntityType: log.EntityType,
		EntityID:   log.EntityID,
		IPAddress:  log.IPAddress,
		CreatedAt:  log.CreatedAt,
	}

	if log.User != nil {
		resp.User = &dto.UserResponse{
			ID:    log.User.ID,
			Email: log.User.Email,
			Name:  log.User.Name,
			Role:  string(log.User.Role),
		}
	}

	if log.Details != nil {
		var details interface{}
		json.Unmarshal(log.Details, &details)
		resp.Details = details
	}

	return resp
}
