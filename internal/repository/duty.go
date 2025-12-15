package repository

import (
	"time"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DutyRepository struct {
	db *gorm.DB
}

func NewDutyRepository(db *gorm.DB) *DutyRepository {
	return &DutyRepository{db: db}
}

func (r *DutyRepository) Create(duty *model.Duty) error {
	return r.db.Create(duty).Error
}

func (r *DutyRepository) CreateBatch(duties []model.Duty) error {
	return r.db.CreateInBatches(duties, 100).Error
}

func (r *DutyRepository) FindByID(id uuid.UUID) (*model.Duty, error) {
	var duty model.Duty
	err := r.db.Preload("Assignee").First(&duty, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &duty, nil
}

func (r *DutyRepository) FindAll(query dto.DutyQuery) ([]model.Duty, int64, error) {
	var duties []model.Duty
	var total int64

	db := r.db.Model(&model.Duty{}).Preload("Assignee")

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.AssigneeID != uuid.Nil {
		db = db.Where("assignee_id = ?", query.AssigneeID)
	}
	if query.StartDate != "" {
		startDate, _ := time.Parse("2006-01-02", query.StartDate)
		db = db.Where("date >= ?", startDate)
	}
	if query.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", query.EndDate)
		db = db.Where("date <= ?", endDate)
	}

	db.Count(&total)

	offset := (query.Page - 1) * query.Limit
	err := db.Offset(offset).Limit(query.Limit).Order("date").Find(&duties).Error

	return duties, total, err
}

func (r *DutyRepository) Update(duty *model.Duty) error {
	return r.db.Save(duty).Error
}

func (r *DutyRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Duty{}, "id = ?", id).Error
}

func (r *DutyRepository) ExistsByDateAndType(date time.Time, dutyType model.DutyType, floor *int) (bool, error) {
	var count int64
	db := r.db.Model(&model.Duty{}).Where("date = ? AND type = ?", date, dutyType)
	if floor != nil {
		db = db.Where("floor = ?", *floor)
	}
	err := db.Count(&count).Error
	return count > 0, err
}
