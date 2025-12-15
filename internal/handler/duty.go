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
	dutyService     *service.DutyService
	swapService     *service.DutySwapRequestService
	auditService    *service.AuditService
}

func NewDutyHandler(dutyService *service.DutyService, swapService *service.DutySwapRequestService, auditService *service.AuditService) *DutyHandler {
	return &DutyHandler{dutyService: dutyService, swapService: swapService, auditService: auditService}
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
// @Success 200 {object} dto.Response{data=[]dto.DutyResponse}
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

	duties, err := h.dutyService.GetAll(query)
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

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
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

// CreateSwapRequest godoc
// @Summary 당직 교대 신청
// @Description 다른 사람의 당직과 교대 신청
// @Tags 당직 교대
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "내 당직 ID"
// @Param request body dto.CreateDutySwapRequest true "교대 대상 당직 ID"
// @Success 201 {object} dto.Response{data=dto.DutySwapRequestResponse}
// @Failure 400 {object} dto.Response
// @Router /duties/{id}/swap-requests [post]
func (h *DutyHandler) CreateSwapRequest(c *gin.Context) {
	sourceDutyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid duty id",
		})
		return
	}

	var req dto.CreateDutySwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	swapReq, err := h.swapService.Create(userID, sourceDutyID, req.TargetDutyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionRequestDutySwap, "duty_swap_request", &swapReq.ID, map[string]interface{}{
		"sourceDutyId": sourceDutyID,
		"targetDutyId": req.TargetDutyID,
	}, c.ClientIP())

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    toSwapRequestResponse(swapReq),
	})
}

// GetPendingSwapRequests godoc
// @Summary 받은 교대 신청 목록
// @Description 나에게 온 대기 중인 교대 신청 목록 조회
// @Tags 당직 교대
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response{data=[]dto.DutySwapRequestResponse}
// @Router /duty-swap-requests/pending [get]
func (h *DutyHandler) GetPendingSwapRequests(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	requests, err := h.swapService.GetPendingForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.DutySwapRequestResponse{}
	for _, r := range requests {
		responses = append(responses, toSwapRequestResponse(&r))
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// GetMySwapRequests godoc
// @Summary 내가 신청한 교대 목록
// @Description 내가 신청한 교대 목록 조회
// @Tags 당직 교대
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response{data=[]dto.DutySwapRequestResponse}
// @Router /duty-swap-requests/my [get]
func (h *DutyHandler) GetMySwapRequests(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	requests, err := h.swapService.GetMyRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.DutySwapRequestResponse{}
	for _, r := range requests {
		responses = append(responses, toSwapRequestResponse(&r))
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// ApproveSwapRequest godoc
// @Summary 교대 신청 승인
// @Description 받은 교대 신청 승인
// @Tags 당직 교대
// @Produce json
// @Security BearerAuth
// @Param id path string true "교대 신청 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /duty-swap-requests/{id}/approve [patch]
func (h *DutyHandler) ApproveSwapRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid swap request id",
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	if err := h.swapService.Approve(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionApproveDutySwap, "duty_swap_request", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

// RejectSwapRequest godoc
// @Summary 교대 신청 거절
// @Description 받은 교대 신청 거절
// @Tags 당직 교대
// @Produce json
// @Security BearerAuth
// @Param id path string true "교대 신청 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /duty-swap-requests/{id}/reject [patch]
func (h *DutyHandler) RejectSwapRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid swap request id",
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	if err := h.swapService.Reject(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(userID, model.AuditActionRejectDutySwap, "duty_swap_request", &id, nil, c.ClientIP())

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

func toSwapRequestResponse(r *model.DutySwapRequest) dto.DutySwapRequestResponse {
	resp := dto.DutySwapRequestResponse{
		ID:        r.ID,
		Status:    string(r.Status),
		CreatedAt: r.CreatedAt,
	}

	if r.Requester != nil {
		resp.Requester = &dto.UserResponse{
			ID:    r.Requester.ID,
			Email: r.Requester.Email,
			Name:  r.Requester.Name,
			Role:  string(r.Requester.Role),
		}
	}

	if r.SourceDuty != nil {
		dutyResp := toDutyResponseWithRelations(r.SourceDuty)
		resp.SourceDuty = &dutyResp
	}

	if r.TargetDuty != nil {
		dutyResp := toDutyResponseWithRelations(r.TargetDuty)
		resp.TargetDuty = &dutyResp
	}

	return resp
}
