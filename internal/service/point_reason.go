package service

import (
	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/repository"

	"github.com/google/uuid"
)

type PointReasonService struct {
	reasonRepo *repository.PointReasonRepository
}

func NewPointReasonService(reasonRepo *repository.PointReasonRepository) *PointReasonService {
	return &PointReasonService{reasonRepo: reasonRepo}
}

func (s *PointReasonService) Create(req dto.CreatePointReasonRequest) (*model.PointReason, error) {
	reason := &model.PointReason{
		Name:  req.Name,
		Type:  model.PointType(req.Type),
		Score: req.Score,
	}

	if err := s.reasonRepo.Create(reason); err != nil {
		return nil, err
	}

	return reason, nil
}

func (s *PointReasonService) GetByID(id uuid.UUID) (*model.PointReason, error) {
	return s.reasonRepo.FindByID(id)
}

func (s *PointReasonService) GetAll() ([]model.PointReason, error) {
	return s.reasonRepo.FindAll()
}

func (s *PointReasonService) GetByType(pointType string) ([]model.PointReason, error) {
	return s.reasonRepo.FindByType(model.PointType(pointType))
}

func (s *PointReasonService) Update(id uuid.UUID, req dto.UpdatePointReasonRequest) (*model.PointReason, error) {
	reason, err := s.reasonRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		reason.Name = req.Name
	}
	if req.Type != "" {
		reason.Type = model.PointType(req.Type)
	}
	if req.Score > 0 {
		reason.Score = req.Score
	}

	if err := s.reasonRepo.Update(reason); err != nil {
		return nil, err
	}

	return reason, nil
}

func (s *PointReasonService) Delete(id uuid.UUID) error {
	return s.reasonRepo.Delete(id)
}
