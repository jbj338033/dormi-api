package model

import (
	"time"

	"github.com/google/uuid"
)

type PointReason struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Type      PointType `gorm:"type:varchar(20);not null"`
	Score     int       `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
