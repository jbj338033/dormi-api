package handler

import (
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService  *service.AuthService
	auditService *service.AuditService
}

func NewAuthHandler(authService *service.AuthService, auditService *service.AuditService) *AuthHandler {
	return &AuthHandler{authService: authService, auditService: auditService}
}

// Login godoc
// @Summary 로그인
// @Description 이메일과 비밀번호로 로그인하여 JWT 토큰 발급
// @Tags 인증
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "로그인 정보"
// @Success 200 {object} dto.Response{data=dto.LoginResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	h.auditService.Log(resp.User.ID, model.AuditActionLogin, "user", &resp.User.ID, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    resp,
	})
}

// CreateUser godoc
// @Summary 사용자 생성
// @Description 새로운 사용자 계정 생성 (관리자 전용)
// @Tags 사용자
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "사용자 정보"
// @Success 201 {object} dto.Response{data=dto.UserResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Router /users [post]
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := h.authService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCreate, "user", &user.ID, map[string]string{"email": user.Email}, c.ClientIP())

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	})
}

// GetAllUsers godoc
// @Summary 사용자 목록
// @Description 모든 사용자 목록 조회 (관리자 전용)
// @Tags 사용자
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response{data=[]dto.UserResponse}
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Router /users [get]
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.UserResponse{}
	for _, u := range users {
		responses = append(responses, dto.UserResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
			Role:  string(u.Role),
		})
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// GetUserByID godoc
// @Summary 사용자 상세
// @Description ID로 사용자 정보 조회 (관리자 전용)
// @Tags 사용자
// @Produce json
// @Security BearerAuth
// @Param id path string true "사용자 ID"
// @Success 200 {object} dto.Response{data=dto.UserResponse}
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid user id",
		})
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Error:   "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	})
}

// UpdateUser godoc
// @Summary 사용자 수정
// @Description 사용자 정보 수정 (관리자 전용)
// @Tags 사용자
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "사용자 ID"
// @Param request body dto.UpdateUserRequest true "수정할 정보"
// @Success 200 {object} dto.Response{data=dto.UserResponse}
// @Failure 400 {object} dto.Response
// @Router /users/{id} [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid user id",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := h.authService.UpdateUser(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionUpdate, "user", &user.ID, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	})
}

// DeleteUser godoc
// @Summary 사용자 삭제
// @Description 사용자 삭제 (관리자 전용)
// @Tags 사용자
// @Produce json
// @Security BearerAuth
// @Param id path string true "사용자 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid user id",
		})
		return
	}

	if err := h.authService.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionDelete, "user", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

// ChangePassword godoc
// @Summary 비밀번호 변경
// @Description 본인 계정 비밀번호 변경
// @Tags 인증
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "비밀번호 변경 정보"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /auth/password [patch]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.authService.ChangePassword(userID, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}
