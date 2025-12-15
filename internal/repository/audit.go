package repository

import (
	"dormi-api/internal/dto"
	"dormi-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditRepository) FindAll(query dto.AuditQuery) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	db := r.db.Model(&model.AuditLog{}).Preload("User")

	if query.UserID != uuid.Nil {
		db = db.Where("user_id = ?", query.UserID)
	}
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}
	if query.EntityType != "" {
		db = db.Where("entity_type = ?", query.EntityType)
	}
	if !query.StartDate.IsZero() {
		db = db.Where("created_at >= ?", query.StartDate)
	}
	if !query.EndDate.IsZero() {
		db = db.Where("created_at <= ?", query.EndDate)
	}

	db.Count(&total)

	offset := (query.Page - 1) * query.Limit
	err := db.Offset(offset).Limit(query.Limit).Order("created_at DESC").Find(&logs).Error

	return logs, total, err
}
