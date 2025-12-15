package model

import (
	"time"

	"github.com/google/uuid"
)

type PointType string

const (
	PointTypeReward  PointType = "REWARD"
	PointTypePenalty PointType = "PENALTY"
)

type Point struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StudentID   uuid.UUID    `gorm:"type:uuid;not null;index"`
	Student     *Student     `gorm:"foreignKey:StudentID"`
	ReasonID    uuid.UUID    `gorm:"type:uuid;not null;index"`
	Reason      *PointReason `gorm:"foreignKey:ReasonID"`
	GivenBy     uuid.UUID    `gorm:"type:uuid;not null"`
	GivenByUser *User        `gorm:"foreignKey:GivenBy"`
	GivenAt     time.Time    `gorm:"not null"`
	Cancelled   bool         `gorm:"default:false"`
	CancelledAt *time.Time
	CancelledBy *uuid.UUID `gorm:"type:uuid"`
}
