package repository

import (
	"dormi-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PointReasonRepository struct {
	db *gorm.DB
}

func NewPointReasonRepository(db *gorm.DB) *PointReasonRepository {
	return &PointReasonRepository{db: db}
}

func (r *PointReasonRepository) Create(reason *model.PointReason) error {
	return r.db.Create(reason).Error
}

func (r *PointReasonRepository) FindByID(id uuid.UUID) (*model.PointReason, error) {
	var reason model.PointReason
	err := r.db.First(&reason, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reason, nil
}

func (r *PointReasonRepository) FindAll() ([]model.PointReason, error) {
	var reasons []model.PointReason
	err := r.db.Order("type, name").Find(&reasons).Error
	return reasons, err
}

func (r *PointReasonRepository) FindByType(pointType model.PointType) ([]model.PointReason, error) {
	var reasons []model.PointReason
	err := r.db.Where("type = ?", pointType).Order("name").Find(&reasons).Error
	return reasons, err
}

func (r *PointReasonRepository) Update(reason *model.PointReason) error {
	return r.db.Save(reason).Error
}

func (r *PointReasonRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.PointReason{}, "id = ?", id).Error
}
