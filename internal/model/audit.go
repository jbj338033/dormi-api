package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditAction string

const (
	AuditActionCreate       AuditAction = "CREATE"
	AuditActionUpdate       AuditAction = "UPDATE"
	AuditActionDelete       AuditAction = "DELETE"
	AuditActionLogin        AuditAction = "LOGIN"
	AuditActionGivePoint    AuditAction = "GIVE_POINT"
	AuditActionCancelPoint  AuditAction = "CANCEL_POINT"
	AuditActionResetPoints  AuditAction = "RESET_POINTS"
	AuditActionSwapDuty     AuditAction = "SWAP_DUTY"
	AuditActionCompleteDuty AuditAction = "COMPLETE_DUTY"
)

type AuditLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	User       *User          `gorm:"foreignKey:UserID"`
	Action     AuditAction    `gorm:"type:varchar(50);not null"`
	EntityType string         `gorm:"type:varchar(50)"`
	EntityID   *uuid.UUID     `gorm:"type:uuid"`
	Details    datatypes.JSON `gorm:"type:jsonb"`
	IPAddress  string         `gorm:"type:varchar(45)"`
	CreatedAt  time.Time      `gorm:"index"`
}
