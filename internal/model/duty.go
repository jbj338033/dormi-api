package model

import (
	"time"

	"github.com/google/uuid"
)

type DutyType string

const (
	DutyTypeDorm       DutyType = "DORM"
	DutyTypeNightStudy DutyType = "NIGHT_STUDY"
)

type Duty struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type       DutyType  `gorm:"type:varchar(20);not null"`
	Date       time.Time `gorm:"type:date;not null;index"`
	Floor      *int      `gorm:"type:int"`
	AssigneeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Assignee   *User     `gorm:"foreignKey:AssigneeID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DutySwapRequestStatus string

const (
	DutySwapStatusPending  DutySwapRequestStatus = "PENDING"
	DutySwapStatusApproved DutySwapRequestStatus = "APPROVED"
	DutySwapStatusRejected DutySwapRequestStatus = "REJECTED"
)

type DutySwapRequest struct {
	ID           uuid.UUID             `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequesterID  uuid.UUID             `gorm:"type:uuid;not null;index"`
	Requester    *User                 `gorm:"foreignKey:RequesterID"`
	SourceDutyID uuid.UUID             `gorm:"type:uuid;not null"`
	SourceDuty   *Duty                 `gorm:"foreignKey:SourceDutyID"`
	TargetDutyID uuid.UUID             `gorm:"type:uuid;not null"`
	TargetDuty   *Duty                 `gorm:"foreignKey:TargetDutyID"`
	Status       DutySwapRequestStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
