package handler

import (
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentHandler struct {
	studentService *service.StudentService
	auditService   *service.AuditService
}

func NewStudentHandler(studentService *service.StudentService, auditService *service.AuditService) *StudentHandler {
	return &StudentHandler{studentService: studentService, auditService: auditService}
}

// Create godoc
// @Summary 학생 생성
// @Description 새로운 학생 등록
// @Tags 학생
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateStudentRequest true "학생 정보"
// @Success 201 {object} dto.Response{data=dto.StudentResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /students [post]
func (h *StudentHandler) Create(c *gin.Context) {
	var req dto.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	student, err := h.studentService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCreate, "student", &student.ID, nil, c.ClientIP())

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    toStudentResponse(student),
	})
}

// GetByID godoc
// @Summary 학생 상세 조회
// @Description ID로 학생 정보 조회
// @Tags 학생
// @Produce json
// @Security BearerAuth
// @Param id path string true "학생 ID"
// @Success 200 {object} dto.Response{data=dto.StudentResponse}
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /students/{id} [get]
func (h *StudentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid student id",
		})
		return
	}

	student, err := h.studentService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Error:   "student not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    toStudentResponse(student),
	})
}

// GetAll godoc
// @Summary 학생 목록 조회
// @Description 학생 목록 조회 (검색, 필터링 지원)
// @Tags 학생
// @Produce json
// @Security BearerAuth
// @Param search query string false "검색어 (이름, 학번)"
// @Param grade query int false "학년"
// @Param room query string false "방 번호"
// @Success 200 {object} dto.Response{data=[]dto.StudentResponse}
// @Failure 400 {object} dto.Response
// @Router /students [get]
func (h *StudentHandler) GetAll(c *gin.Context) {
	var query dto.StudentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	students, err := h.studentService.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	responses := []dto.StudentResponse{}
	for _, s := range students {
		responses = append(responses, toStudentResponse(&s))
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    responses,
	})
}

// Update godoc
// @Summary 학생 수정
// @Description 학생 정보 수정
// @Tags 학생
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "학생 ID"
// @Param request body dto.UpdateStudentRequest true "수정할 정보"
// @Success 200 {object} dto.Response{data=dto.StudentResponse}
// @Failure 400 {object} dto.Response
// @Router /students/{id} [put]
func (h *StudentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid student id",
		})
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	student, err := h.studentService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionUpdate, "student", &student.ID, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    toStudentResponse(student),
	})
}

// Delete godoc
// @Summary 학생 삭제
// @Description 학생 정보 삭제 (soft delete)
// @Tags 학생
// @Produce json
// @Security BearerAuth
// @Param id path string true "학생 ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /students/{id} [delete]
func (h *StudentHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "invalid student id",
		})
		return
	}

	if err := h.studentService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionDelete, "student", &id, nil, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
	})
}

// Import godoc
// @Summary CSV 일괄 등록
// @Description CSV 파일로 학생 일괄 등록
// @Tags 학생
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CSV 파일"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /students/import [post]
func (h *StudentHandler) Import(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   "file is required",
		})
		return
	}
	defer file.Close()

	count, err := h.studentService.ImportCSV(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	h.auditService.Log(userID, model.AuditActionCreate, "student", nil, map[string]int{"count": count}, c.ClientIP())

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    map[string]int{"imported": count},
	})
}

func toStudentResponse(s *model.Student) dto.StudentResponse {
	return dto.StudentResponse{
		ID:            s.ID,
		StudentNumber: s.StudentNumber,
		Name:          s.Name,
		RoomNumber:    s.RoomNumber,
		Grade:         s.Grade,
		CreatedAt:     s.CreatedAt,
	}
}
