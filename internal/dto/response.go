package dto

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Pagination `json:"meta"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}

type StudentResponse struct {
	ID            uuid.UUID `json:"id"`
	StudentNumber string    `json:"studentNumber"`
	Name          string    `json:"name"`
	RoomNumber    string    `json:"roomNumber"`
	Grade         int       `json:"grade"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PointReasonResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Score int       `json:"score"`
}

type PointResponse struct {
	ID          uuid.UUID            `json:"id"`
	Student     *StudentResponse     `json:"student,omitempty"`
	Reason      *PointReasonResponse `json:"reason,omitempty"`
	GivenBy     *UserResponse        `json:"givenBy,omitempty"`
	GivenAt     time.Time            `json:"givenAt"`
	Cancelled   bool                 `json:"cancelled"`
	CancelledAt *time.Time           `json:"cancelledAt,omitempty"`
}

type PointSummary struct {
	StudentID    uuid.UUID `json:"studentId"`
	TotalReward  int       `json:"totalReward"`
	TotalPenalty int       `json:"totalPenalty"`
	NetScore     int       `json:"netScore"`
}

type DutyResponse struct {
	ID        uuid.UUID     `json:"id"`
	Type      string        `json:"type"`
	Date      string        `json:"date"`
	Floor     *int          `json:"floor,omitempty"`
	Assignee  *UserResponse `json:"assignee,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
}

type DutySwapRequestResponse struct {
	ID         uuid.UUID     `json:"id"`
	Requester  *UserResponse `json:"requester,omitempty"`
	SourceDuty *DutyResponse `json:"sourceDuty,omitempty"`
	TargetDuty *DutyResponse `json:"targetDuty,omitempty"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"createdAt"`
}

type AuditLogResponse struct {
	ID         uuid.UUID     `json:"id"`
	User       *UserResponse `json:"user,omitempty"`
	Action     string        `json:"action"`
	EntityType string        `json:"entityType"`
	EntityID   *uuid.UUID    `json:"entityId,omitempty"`
	Details    interface{}   `json:"details,omitempty"`
	IPAddress  string        `json:"ipAddress"`
	CreatedAt  time.Time     `json:"createdAt"`
}
