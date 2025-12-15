package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin      Role = "ADMIN"
	RoleSupervisor Role = "SUPERVISOR"
	RoleCouncil    Role = "COUNCIL"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Role      Role      `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
