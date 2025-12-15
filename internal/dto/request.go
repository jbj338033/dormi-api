package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=ADMIN SUPERVISOR COUNCIL"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Name     string `json:"name"`
	Role     string `json:"role" binding:"omitempty,oneof=ADMIN SUPERVISOR COUNCIL"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}

type CreateStudentRequest struct {
	StudentNumber string `json:"studentNumber" binding:"required"`
	Name          string `json:"name" binding:"required"`
	RoomNumber    string `json:"roomNumber" binding:"required"`
	Grade         int    `json:"grade" binding:"required,min=1,max=3"`
}

type UpdateStudentRequest struct {
	StudentNumber string `json:"studentNumber"`
	Name          string `json:"name"`
	RoomNumber    string `json:"roomNumber"`
	Grade         int    `json:"grade" binding:"omitempty,min=1,max=3"`
}

type StudentQuery struct {
	Search string `form:"search"`
	Grade  int    `form:"grade"`
	Room   string `form:"room"`
}

type GivePointRequest struct {
	StudentID uuid.UUID `json:"studentId" binding:"required"`
	ReasonID  uuid.UUID `json:"reasonId" binding:"required"`
}

type BulkGivePointRequest struct {
	StudentIDs []uuid.UUID `json:"studentIds" binding:"required,min=1"`
	ReasonID   uuid.UUID   `json:"reasonId" binding:"required"`
}

type CreatePointReasonRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=REWARD PENALTY"`
	Score int    `json:"score" binding:"required,min=1"`
}

type UpdatePointReasonRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type" binding:"omitempty,oneof=REWARD PENALTY"`
	Score int    `json:"score" binding:"omitempty,min=1"`
}

type PointQuery struct {
	StudentID uuid.UUID `form:"studentId"`
	Type      string    `form:"type"`
	StartDate string    `form:"startDate"`
	EndDate   string    `form:"endDate"`
}

type CreateDutyRequest struct {
	Type       string    `json:"type" binding:"required,oneof=DORM NIGHT_STUDY"`
	Date       string    `json:"date" binding:"required"`
	Floor      *int      `json:"floor"`
	AssigneeID uuid.UUID `json:"assigneeId" binding:"required"`
}

type UpdateDutyRequest struct {
	Type       string    `json:"type" binding:"omitempty,oneof=DORM NIGHT_STUDY"`
	Date       string    `json:"date"`
	Floor      *int      `json:"floor"`
	AssigneeID uuid.UUID `json:"assigneeId"`
}

type GenerateDutyRequest struct {
	Type        string      `json:"type" binding:"required,oneof=DORM NIGHT_STUDY"`
	StartDate   string      `json:"startDate" binding:"required"`
	EndDate     string      `json:"endDate" binding:"required"`
	AssigneeIDs []uuid.UUID `json:"assigneeIds" binding:"required,min=1"`
	Floor       *int        `json:"floor"`
}

type CreateDutySwapRequest struct {
	TargetDutyID uuid.UUID `json:"targetDutyId" binding:"required"`
}

type DutyQuery struct {
	Type       string    `form:"type"`
	AssigneeID uuid.UUID `form:"assigneeId"`
	StartDate  string    `form:"startDate"`
	EndDate    string    `form:"endDate"`
}

type AuditQuery struct {
	UserID     uuid.UUID `form:"userId"`
	Action     string    `form:"action"`
	EntityType string    `form:"entityType"`
	StartDate  time.Time `form:"startDate"`
	EndDate    time.Time `form:"endDate"`
	Page       int       `form:"page,default=1"`
	Limit      int       `form:"limit,default=20"`
}
