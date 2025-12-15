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

func (r *DutyRepository) FindAll(query dto.DutyQuery) ([]model.Duty, error) {
	var duties []model.Duty

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

	err := db.Order("date").Find(&duties).Error

	return duties, err
}

func (r *DutyRepository) Update(duty *model.Duty) error {
	return r.db.Save(duty).Error
}

func (r *DutyRepository) UpdateAssignee(id, assigneeID uuid.UUID) error {
	return r.db.Model(&model.Duty{}).Where("id = ?", id).Update("assignee_id", assigneeID).Error
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

type DutySwapRequestRepository struct {
	db *gorm.DB
}

func NewDutySwapRequestRepository(db *gorm.DB) *DutySwapRequestRepository {
	return &DutySwapRequestRepository{db: db}
}

func (r *DutySwapRequestRepository) Create(req *model.DutySwapRequest) error {
	return r.db.Create(req).Error
}

func (r *DutySwapRequestRepository) FindByID(id uuid.UUID) (*model.DutySwapRequest, error) {
	var req model.DutySwapRequest
	err := r.db.Preload("Requester").Preload("SourceDuty").Preload("SourceDuty.Assignee").Preload("TargetDuty").Preload("TargetDuty.Assignee").First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *DutySwapRequestRepository) FindPendingByTargetAssignee(assigneeID uuid.UUID) ([]model.DutySwapRequest, error) {
	var requests []model.DutySwapRequest
	err := r.db.
		Preload("Requester").
		Preload("SourceDuty").
		Preload("SourceDuty.Assignee").
		Preload("TargetDuty").
		Preload("TargetDuty.Assignee").
		Joins("JOIN duties ON duties.id = duty_swap_requests.target_duty_id").
		Where("duties.assignee_id = ? AND duty_swap_requests.status = ?", assigneeID, model.DutySwapStatusPending).
		Order("duty_swap_requests.created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *DutySwapRequestRepository) FindByRequester(requesterID uuid.UUID) ([]model.DutySwapRequest, error) {
	var requests []model.DutySwapRequest
	err := r.db.
		Preload("Requester").
		Preload("SourceDuty").
		Preload("SourceDuty.Assignee").
		Preload("TargetDuty").
		Preload("TargetDuty.Assignee").
		Where("requester_id = ?", requesterID).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *DutySwapRequestRepository) Update(req *model.DutySwapRequest) error {
	return r.db.Save(req).Error
}

func (r *DutySwapRequestRepository) ExistsPendingBetweenDuties(sourceDutyID, targetDutyID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.DutySwapRequest{}).
		Where("source_duty_id = ? AND target_duty_id = ? AND status = ?", sourceDutyID, targetDutyID, model.DutySwapStatusPending).
		Count(&count).Error
	return count > 0, err
}
