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
	Completed  bool      `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
