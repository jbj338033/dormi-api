package repository

import (
	"time"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PointRepository struct {
	db *gorm.DB
}

func NewPointRepository(db *gorm.DB) *PointRepository {
	return &PointRepository{db: db}
}

func (r *PointRepository) Create(point *model.Point) error {
	return r.db.Create(point).Error
}

func (r *PointRepository) CreateBatch(points []model.Point) error {
	return r.db.CreateInBatches(points, 100).Error
}

func (r *PointRepository) FindByID(id uuid.UUID) (*model.Point, error) {
	var point model.Point
	err := r.db.Preload("Student").Preload("Reason").Preload("GivenByUser").First(&point, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &point, nil
}

func (r *PointRepository) FindAll(query dto.PointQuery) ([]model.Point, error) {
	var points []model.Point

	db := r.db.Model(&model.Point{}).Preload("Student").Preload("Reason").Preload("GivenByUser")

	if query.StudentID != uuid.Nil {
		db = db.Where("student_id = ?", query.StudentID)
	}
	if query.Type != "" {
		db = db.Joins("JOIN point_reasons ON point_reasons.id = points.reason_id").
			Where("point_reasons.type = ?", query.Type)
	}
	if query.StartDate != "" {
		db = db.Where("given_at >= ?", query.StartDate)
	}
	if query.EndDate != "" {
		db = db.Where("given_at <= ?", query.EndDate)
	}

	err := db.Order("given_at DESC").Find(&points).Error

	return points, err
}

func (r *PointRepository) FindByStudentID(studentID uuid.UUID) ([]model.Point, error) {
	var points []model.Point
	err := r.db.Preload("Reason").Preload("GivenByUser").
		Where("student_id = ? AND cancelled = false", studentID).
		Order("given_at DESC").
		Find(&points).Error
	return points, err
}

func (r *PointRepository) GetSummary(studentID uuid.UUID) (*dto.PointSummary, error) {
	var result dto.PointSummary
	result.StudentID = studentID

	err := r.db.Model(&model.Point{}).
		Joins("JOIN point_reasons ON point_reasons.id = points.reason_id").
		Select("COALESCE(SUM(CASE WHEN point_reasons.type = 'REWARD' THEN point_reasons.score ELSE 0 END), 0) as total_reward, COALESCE(SUM(CASE WHEN point_reasons.type = 'PENALTY' THEN point_reasons.score ELSE 0 END), 0) as total_penalty").
		Where("points.student_id = ? AND points.cancelled = false", studentID).
		Row().Scan(&result.TotalReward, &result.TotalPenalty)

	result.NetScore = result.TotalReward - result.TotalPenalty
	return &result, err
}

func (r *PointRepository) Cancel(id uuid.UUID, cancelledBy uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.Point{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"cancelled":    true,
			"cancelled_at": now,
			"cancelled_by": cancelledBy,
		}).Error
}

func (r *PointRepository) ResetAll() error {
	return r.db.Where("1 = 1").Delete(&model.Point{}).Error
}
