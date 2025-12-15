package service

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/repository"

	"github.com/google/uuid"
)

type StudentService struct {
	studentRepo *repository.StudentRepository
}

func NewStudentService(studentRepo *repository.StudentRepository) *StudentService {
	return &StudentService{studentRepo: studentRepo}
}

func (s *StudentService) Create(req dto.CreateStudentRequest) (*model.Student, error) {
	exists, err := s.studentRepo.ExistsByStudentNumber(req.StudentNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("student number already exists")
	}

	student := &model.Student{
		StudentNumber: req.StudentNumber,
		Name:          req.Name,
		RoomNumber:    req.RoomNumber,
		Grade:         req.Grade,
	}

	if err := s.studentRepo.Create(student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *StudentService) GetByID(id uuid.UUID) (*model.Student, error) {
	return s.studentRepo.FindByID(id)
}

func (s *StudentService) GetAll(query dto.StudentQuery) ([]model.Student, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	return s.studentRepo.FindAll(query)
}

func (s *StudentService) Update(id uuid.UUID, req dto.UpdateStudentRequest) (*model.Student, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.StudentNumber != "" && req.StudentNumber != student.StudentNumber {
		exists, err := s.studentRepo.ExistsByStudentNumber(req.StudentNumber)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("student number already exists")
		}
		student.StudentNumber = req.StudentNumber
	}
	if req.Name != "" {
		student.Name = req.Name
	}
	if req.RoomNumber != "" {
		student.RoomNumber = req.RoomNumber
	}
	if req.Grade > 0 {
		student.Grade = req.Grade
	}

	if err := s.studentRepo.Update(student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *StudentService) Delete(id uuid.UUID) error {
	return s.studentRepo.Delete(id)
}

func (s *StudentService) ImportCSV(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)

	_, err := csvReader.Read()
	if err != nil {
		return 0, err
	}

	var students []model.Student
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		if len(record) < 4 {
			continue
		}

		grade, err := strconv.Atoi(record[3])
		if err != nil {
			continue
		}

		students = append(students, model.Student{
			StudentNumber: record[0],
			Name:          record[1],
			RoomNumber:    record[2],
			Grade:         grade,
		})
	}

	if len(students) == 0 {
		return 0, errors.New("no valid records found")
	}

	if err := s.studentRepo.CreateBatch(students); err != nil {
		return 0, err
	}

	return len(students), nil
}
