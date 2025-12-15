package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StudentNumber string         `gorm:"type:varchar(20);uniqueIndex;not null"`
	Name          string         `gorm:"type:varchar(100);not null"`
	RoomNumber    string         `gorm:"type:varchar(20);not null"`
	Grade         int            `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
