package handler

import (
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PointReasonHandler struct {
	reasonService *service.PointReasonService
	auditService  *service.AuditService
}

func NewPointReasonHandler(reasonService *service.PointReasonService, auditService *service.AuditService) *PointReasonHandler {
	return &PointReasonHandler{reasonService: reasonService, auditService: auditService}
}

// Create godoc
// @Summary 상벌점 사유 생성
// @Description 새로운 상벌점 사유 등록
// @Tags 상벌점사유
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePointReasonRequest true "사유 정보"
// @Success 201 {object} dto.Response{data=dto.PointReasonResponse}
// @Failure 400 {object} dto.Response
// @Router /point-reasons [post]
func (h *PointReasonHandler) Create(c *gin.Context) {
	var req dto.CreatePointReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	reason, err := h.reasonService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCreate, "point_reason", &reason.ID, nil, c.ClientIP())

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    toPointReasonResponse(reason),
	})
}

// GetByID godoc
// @Summary 상벌점 사유 상세
// @Description ID로 상벌점 사유 조회
// @Tags 상벌점사유
// @Produce json
// @Security BearerAuth
// @Param id path string true "사유 ID"
// @Success 200 {object} dto.Response{data=dto.PointReasonResponse}
// @Failure 404 {object} dto.Response
// @Router /point-reasons/{id} [get]
func (h *PointReasonHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid reason id",
		})
		return
	}

	reason, err := h.reasonService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Error:   "reason not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    toPointReasonResponse(reason),
	})
}

// GetAll godoc
// @Summary 상벌점 사유 목록
// @Description 모든 상벌점 사유 조회
// @Tags 상벌점사유
// @Produce json
// @Security BearerAuth
// @Param type query string false "유형 (REWARD, PENALTY)"
// @Success 200 {object} dto.Response{data=[]dto.PointReasonResponse}
// @Router /point-reasons [get]
func (h *PointReasonHandler) GetAll(c *gin.Context) {
	pointType := c.Query("type")

	var reasons []model.PointReason
	var err error

	if pointType != "" {
		reasons, err = h.reasonService.GetByType(pointType)
	} else {
		reasons, err = h.reasonService.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.PointReasonResponse{}
	for _, r := range reasons {
		responses = append(responses, toPointReasonResponse(&r))
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// Update godoc
// @Summary 상벌점 사유 수정
// @Description 상벌점 사유 수정
// @Tags 상벌점사유
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "사유 ID"
// @Param request body dto.UpdatePointReasonRequest true "수정할 정보"
// @Success 200 {object} dto.Response{data=dto.PointReasonResponse}
// @Failure 400 {object} dto.Response
// @Router /point-reasons/{id} [put]
func (h *PointReasonHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid reason id",
		})
		return
	}

	var req dto.UpdatePointReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	reason, err := h.reasonService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionUpdate, "point_reason", &reason.ID, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    toPointReasonResponse(reason),
	})
}

// Delete godoc
// @Summary 상벌점 사유 삭제
// @Description 상벌점 사유 삭제
// @Tags 상벌점사유
// @Produce json
// @Security BearerAuth
// @Param id path string true "사유 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /point-reasons/{id} [delete]
func (h *PointReasonHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid reason id",
		})
		return
	}

	if err := h.reasonService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionDelete, "point_reason", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

func toPointReasonResponse(r *model.PointReason) dto.PointReasonResponse {
	return dto.PointReasonResponse{
		ID:    r.ID,
		Name:  r.Name,
		Type:  string(r.Type),
		Score: r.Score,
	}
}
