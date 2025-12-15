package handler

import (
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PointHandler struct {
	pointService *service.PointService
	auditService *service.AuditService
}

func NewPointHandler(pointService *service.PointService, auditService *service.AuditService) *PointHandler {
	return &PointHandler{pointService: pointService, auditService: auditService}
}

// GivePoint godoc
// @Summary 상벌점 부여
// @Description 학생에게 상벌점 부여
// @Tags 상벌점
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GivePointRequest true "상벌점 정보"
// @Success 201 {object} dto.Response{data=dto.PointResponse}
// @Failure 400 {object} dto.Response
// @Router /points [post]
func (h *PointHandler) GivePoint(c *gin.Context) {
	var req dto.GivePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	point, err := h.pointService.GivePoint(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionGivePoint, "point", &point.ID, map[string]any{
		"studentId": req.StudentID,
		"reasonId":  req.ReasonID,
	}, c.ClientIP())

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    toPointResponse(point),
	})
}

// BulkGivePoints godoc
// @Summary 상벌점 다건 부여
// @Description 여러 학생에게 동시에 상벌점 부여
// @Tags 상벌점
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BulkGivePointRequest true "다건 상벌점 정보"
// @Success 201 {object} dto.Response{data=[]dto.PointResponse}
// @Failure 400 {object} dto.Response
// @Router /points/bulk [post]
func (h *PointHandler) BulkGivePoints(c *gin.Context) {
	var req dto.BulkGivePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	points, err := h.pointService.BulkGivePoints(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionGivePoint, "point", nil, map[string]any{
		"studentIds": req.StudentIDs,
		"reasonId":   req.ReasonID,
		"count":      len(points),
	}, c.ClientIP())

	responses := []dto.PointResponse{}
	for _, p := range points {
		responses = append(responses, toPointResponse(&p))
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// GetAll godoc
// @Summary 상벌점 목록 조회
// @Description 상벌점 목록 조회 (필터링 지원)
// @Tags 상벌점
// @Produce json
// @Security BearerAuth
// @Param studentId query string false "학생 ID"
// @Param type query string false "유형 (REWARD, PENALTY)"
// @Param startDate query string false "시작일 (YYYY-MM-DD)"
// @Param endDate query string false "종료일 (YYYY-MM-DD)"
// @Param page query int false "페이지" default(1)
// @Param limit query int false "페이지당 개수" default(20)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.PointResponse}
// @Failure 400 {object} dto.Response
// @Router /points [get]
func (h *PointHandler) GetAll(c *gin.Context) {
	var query dto.PointQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	points, total, err := h.pointService.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.PointResponse{}
	for _, p := range points {
		responses = append(responses, toPointResponse(&p))
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

// GetByStudentID godoc
// @Summary 학생별 상벌점 조회
// @Description 특정 학생의 상벌점 목록 조회
// @Tags 상벌점
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "학생 ID"
// @Success 200 {object} dto.Response{data=[]dto.PointResponse}
// @Failure 400 {object} dto.Response
// @Router /points/student/{studentId} [get]
func (h *PointHandler) GetByStudentID(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid student id",
		})
		return
	}

	points, err := h.pointService.GetByStudentID(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.PointResponse{}
	for _, p := range points {
		responses = append(responses, toPointResponse(&p))
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// GetSummary godoc
// @Summary 학생별 상벌점 요약
// @Description 특정 학생의 상벌점 요약 정보 조회
// @Tags 상벌점
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "학생 ID"
// @Success 200 {object} dto.Response{data=dto.PointSummary}
// @Failure 400 {object} dto.Response
// @Router /points/student/{studentId}/summary [get]
func (h *PointHandler) GetSummary(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid student id",
		})
		return
	}

	summary, err := h.pointService.GetSummary(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    summary,
	})
}

// Cancel godoc
// @Summary 상벌점 취소
// @Description 상벌점 취소 (soft delete)
// @Tags 상벌점
// @Produce json
// @Security BearerAuth
// @Param id path string true "상벌점 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /points/{id}/cancel [patch]
func (h *PointHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid point id",
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.pointService.Cancel(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionCancelPoint, "point", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

// Reset godoc
// @Summary 상벌점 전체 초기화
// @Description 모든 상벌점 삭제 (관리자 전용)
// @Tags 상벌점
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /points/reset [delete]
func (h *PointHandler) Reset(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.pointService.ResetAll(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionResetPoints, "point", nil, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

func toPointResponse(p *model.Point) dto.PointResponse {
	resp := dto.PointResponse{
		ID:          p.ID,
		GivenAt:     p.GivenAt,
		Cancelled:   p.Cancelled,
		CancelledAt: p.CancelledAt,
	}

	if p.Student != nil {
		student := toStudentResponse(p.Student)
		resp.Student = &student
	}

	if p.Reason != nil {
		resp.Reason = &dto.PointReasonResponse{
			ID:    p.Reason.ID,
			Name:  p.Reason.Name,
			Type:  string(p.Reason.Type),
			Score: p.Reason.Score,
		}
	}

	if p.GivenByUser != nil {
		resp.GivenBy = &dto.UserResponse{
			ID:    p.GivenByUser.ID,
			Email: p.GivenByUser.Email,
			Name:  p.GivenByUser.Name,
			Role:  string(p.GivenByUser.Role),
		}
	}

	return resp
}
