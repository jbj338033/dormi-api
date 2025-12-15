package handler

import (
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DutyHandler struct {
	dutyService  *service.DutyService
	auditService *service.AuditService
}

func NewDutyHandler(dutyService *service.DutyService, auditService *service.AuditService) *DutyHandler {
	return &DutyHandler{dutyService: dutyService, auditService: auditService}
}

// Create godoc
// @Summary 당직 생성
// @Description 새로운 당직 일정 생성
// @Tags 당직
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateDutyRequest true "당직 정보"
// @Success 201 {object} dto.Response{data=dto.DutyResponse}
// @Failure 400 {object} dto.Response
// @Router /duties [post]
func (h *DutyHandler) Create(c *gin.Context) {
	var req dto.CreateDutyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	duty, err := h.dutyService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCreate, "duty", &duty.ID, nil, c.ClientIP())

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    toDutyResponse(duty),
	})
}

// GetByID godoc
// @Summary 당직 상세 조회
// @Description ID로 당직 정보 조회
// @Tags 당직
// @Produce json
// @Security BearerAuth
// @Param id path string true "당직 ID"
// @Success 200 {object} dto.Response{data=dto.DutyResponse}
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /duties/{id} [get]
func (h *DutyHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid duty id",
		})
		return
	}

	duty, err := h.dutyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Error:   "duty not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    toDutyResponseWithRelations(duty),
	})
}

// GetAll godoc
// @Summary 당직 목록 조회
// @Description 당직 목록 조회 (필터링 지원)
// @Tags 당직
// @Produce json
// @Security BearerAuth
// @Param type query string false "유형 (DORM, NIGHT_STUDY)"
// @Param assigneeId query string false "담당자 ID"
// @Param startDate query string false "시작일 (YYYY-MM-DD)"
// @Param endDate query string false "종료일 (YYYY-MM-DD)"
// @Param page query int false "페이지" default(1)
// @Param limit query int false "페이지당 개수" default(20)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.DutyResponse}
// @Failure 400 {object} dto.Response
// @Router /duties [get]
func (h *DutyHandler) GetAll(c *gin.Context) {
	var query dto.DutyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	duties, total, err := h.dutyService.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.DutyResponse{}
	for _, d := range duties {
		responses = append(responses, toDutyResponseWithRelations(&d))
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

// Update godoc
// @Summary 당직 수정
// @Description 당직 정보 수정
// @Tags 당직
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "당직 ID"
// @Param request body dto.UpdateDutyRequest true "수정할 정보"
// @Success 200 {object} dto.Response{data=dto.DutyResponse}
// @Failure 400 {object} dto.Response
// @Router /duties/{id} [put]
func (h *DutyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid duty id",
		})
		return
	}

	var req dto.UpdateDutyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	duty, err := h.dutyService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionUpdate, "duty", &duty.ID, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    toDutyResponse(duty),
	})
}

// Delete godoc
// @Summary 당직 삭제
// @Description 당직 정보 삭제
// @Tags 당직
// @Produce json
// @Security BearerAuth
// @Param id path string true "당직 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /duties/{id} [delete]
func (h *DutyHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid duty id",
		})
		return
	}

	if err := h.dutyService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionDelete, "duty", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

// Generate godoc
// @Summary 당직 자동 생성
// @Description 기간과 담당자 목록으로 당직 자동 생성
// @Tags 당직
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GenerateDutyRequest true "자동 생성 정보"
// @Success 201 {object} dto.Response{data=[]dto.DutyResponse}
// @Failure 400 {object} dto.Response
// @Router /duties/generate [post]
func (h *DutyHandler) Generate(c *gin.Context) {
	var req dto.GenerateDutyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	duties, err := h.dutyService.Generate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCreate, "duty", nil, map[string]int{"count": len(duties)}, c.ClientIP())

	responses := []dto.DutyResponse{}
	for _, d := range duties {
		responses = append(responses, toDutyResponse(&d))
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// Swap godoc
// @Summary 당직 교대
// @Description 두 당직의 담당자 교대
// @Tags 당직
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "당직 ID"
// @Param request body dto.SwapDutyRequest true "교대 대상 당직 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /duties/{id}/swap [post]
func (h *DutyHandler) Swap(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid duty id",
		})
		return
	}

	var req dto.SwapDutyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.dutyService.Swap(id, req.TargetDutyID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionSwapDuty, "duty", &id, map[string]interface{}{
		"targetDutyId": req.TargetDutyID,
	}, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

// Complete godoc
// @Summary 당직 완료 처리
// @Description 당직을 완료 상태로 변경
// @Tags 당직
// @Produce json
// @Security BearerAuth
// @Param id path string true "당직 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /duties/{id}/complete [patch]
func (h *DutyHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid duty id",
		})
		return
	}

	if err := h.dutyService.Complete(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCompleteDuty, "duty", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

func toDutyResponse(d *model.Duty) dto.DutyResponse {
	return dto.DutyResponse{
		ID:        d.ID,
		Type:      string(d.Type),
		Date:      d.Date.Format("2006-01-02"),
		Floor:     d.Floor,
		Completed: d.Completed,
		CreatedAt: d.CreatedAt,
	}
}

func toDutyResponseWithRelations(d *model.Duty) dto.DutyResponse {
	resp := toDutyResponse(d)

	if d.Assignee != nil {
		resp.Assignee = &dto.UserResponse{
			ID:    d.Assignee.ID,
			Email: d.Assignee.Email,
			Name:  d.Assignee.Name,
			Role:  string(d.Assignee.Role),
		}
	}

	return resp
}
