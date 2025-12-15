package service

import (
	"errors"
	"time"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/repository"

	"github.com/google/uuid"
)

type DutyService struct {
	dutyRepo *repository.DutyRepository
}

func NewDutyService(dutyRepo *repository.DutyRepository) *DutyService {
	return &DutyService{dutyRepo: dutyRepo}
}

func (s *DutyService) Create(req dto.CreateDutyRequest) (*model.Duty, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	dutyType := model.DutyType(req.Type)

	weekday := date.Weekday()
	if dutyType == model.DutyTypeDorm {
		if weekday == time.Friday || weekday == time.Saturday {
			return nil, errors.New("DORM duty is only for Sunday to Thursday")
		}
	} else if dutyType == model.DutyTypeNightStudy {
		if weekday < time.Monday || weekday > time.Thursday {
			return nil, errors.New("NIGHT_STUDY duty is only for Monday to Thursday")
		}
		if req.Floor == nil {
			return nil, errors.New("floor is required for NIGHT_STUDY duty")
		}
	}

	duty := &model.Duty{
		Type:       dutyType,
		Date:       date,
		Floor:      req.Floor,
		AssigneeID: req.AssigneeID,
	}

	if err := s.dutyRepo.Create(duty); err != nil {
		return nil, err
	}

	return duty, nil
}

func (s *DutyService) GetByID(id uuid.UUID) (*model.Duty, error) {
	return s.dutyRepo.FindByID(id)
}

func (s *DutyService) GetAll(query dto.DutyQuery) ([]model.Duty, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	return s.dutyRepo.FindAll(query)
}

func (s *DutyService) Update(id uuid.UUID, req dto.UpdateDutyRequest) (*model.Duty, error) {
	duty, err := s.dutyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Type != "" {
		duty.Type = model.DutyType(req.Type)
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		duty.Date = date
	}
	if req.Floor != nil {
		duty.Floor = req.Floor
	}
	if req.AssigneeID != uuid.Nil {
		duty.AssigneeID = req.AssigneeID
	}

	if err := s.dutyRepo.Update(duty); err != nil {
		return nil, err
	}

	return duty, nil
}

func (s *DutyService) Delete(id uuid.UUID) error {
	return s.dutyRepo.Delete(id)
}

func (s *DutyService) Generate(req dto.GenerateDutyRequest) ([]model.Duty, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	dutyType := model.DutyType(req.Type)
	var duties []model.Duty
	assigneeIdx := 0

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		weekday := date.Weekday()

		valid := false
		if dutyType == model.DutyTypeDorm {
			valid = weekday != time.Friday && weekday != time.Saturday
		} else if dutyType == model.DutyTypeNightStudy {
			valid = weekday >= time.Monday && weekday <= time.Thursday
		}

		if valid {
			duty := model.Duty{
				Type:       dutyType,
				Date:       date,
				Floor:      req.Floor,
				AssigneeID: req.AssigneeIDs[assigneeIdx%len(req.AssigneeIDs)],
			}
			duties = append(duties, duty)
			assigneeIdx++
		}
	}

	if len(duties) == 0 {
		return nil, errors.New("no valid duty dates found")
	}

	if err := s.dutyRepo.CreateBatch(duties); err != nil {
		return nil, err
	}

	return duties, nil
}

func (s *DutyService) Swap(id uuid.UUID, targetID uuid.UUID) error {
	duty1, err := s.dutyRepo.FindByID(id)
	if err != nil {
		return errors.New("duty not found")
	}

	duty2, err := s.dutyRepo.FindByID(targetID)
	if err != nil {
		return errors.New("target duty not found")
	}

	if duty1.Type != duty2.Type {
		return errors.New("can only swap duties of the same type")
	}

	duty1.AssigneeID, duty2.AssigneeID = duty2.AssigneeID, duty1.AssigneeID

	if err := s.dutyRepo.Update(duty1); err != nil {
		return err
	}
	if err := s.dutyRepo.Update(duty2); err != nil {
		return err
	}

	return nil
}

func (s *DutyService) Complete(id uuid.UUID) error {
	duty, err := s.dutyRepo.FindByID(id)
	if err != nil {
		return errors.New("duty not found")
	}

	duty.Completed = true
	return s.dutyRepo.Update(duty)
}
