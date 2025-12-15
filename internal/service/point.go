package service

import (
	"errors"
	"time"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/repository"

	"github.com/google/uuid"
)

type PointService struct {
	pointRepo   *repository.PointRepository
	studentRepo *repository.StudentRepository
	reasonRepo  *repository.PointReasonRepository
}

func NewPointService(pointRepo *repository.PointRepository, studentRepo *repository.StudentRepository, reasonRepo *repository.PointReasonRepository) *PointService {
	return &PointService{pointRepo: pointRepo, studentRepo: studentRepo, reasonRepo: reasonRepo}
}

func (s *PointService) GivePoint(req dto.GivePointRequest, givenBy uuid.UUID) (*model.Point, error) {
	_, err := s.reasonRepo.FindByID(req.ReasonID)
	if err != nil {
		return nil, errors.New("invalid reason")
	}

	_, err = s.studentRepo.FindByID(req.StudentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	point := &model.Point{
		StudentID: req.StudentID,
		ReasonID:  req.ReasonID,
		GivenBy:   givenBy,
		GivenAt:   time.Now(),
	}

	if err := s.pointRepo.Create(point); err != nil {
		return nil, err
	}

	return s.pointRepo.FindByID(point.ID)
}

func (s *PointService) BulkGivePoints(req dto.BulkGivePointRequest, givenBy uuid.UUID) ([]model.Point, error) {
	_, err := s.reasonRepo.FindByID(req.ReasonID)
	if err != nil {
		return nil, errors.New("invalid reason")
	}

	var points []model.Point
	now := time.Now()

	for _, studentID := range req.StudentIDs {
		points = append(points, model.Point{
			StudentID: studentID,
			ReasonID:  req.ReasonID,
			GivenBy:   givenBy,
			GivenAt:   now,
		})
	}

	if err := s.pointRepo.CreateBatch(points); err != nil {
		return nil, err
	}

	return points, nil
}

func (s *PointService) GetAll(query dto.PointQuery) ([]model.Point, error) {
	return s.pointRepo.FindAll(query)
}

func (s *PointService) GetByStudentID(studentID uuid.UUID) ([]model.Point, error) {
	return s.pointRepo.FindByStudentID(studentID)
}

func (s *PointService) GetSummary(studentID uuid.UUID) (*dto.PointSummary, error) {
	return s.pointRepo.GetSummary(studentID)
}

func (s *PointService) Cancel(id uuid.UUID, cancelledBy uuid.UUID) error {
	point, err := s.pointRepo.FindByID(id)
	if err != nil {
		return errors.New("point not found")
	}

	if point.Cancelled {
		return errors.New("point already cancelled")
	}

	return s.pointRepo.Cancel(id, cancelledBy)
}

func (s *PointService) ResetAll() error {
	return s.pointRepo.ResetAll()
}
