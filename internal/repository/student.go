package repository

import (
	"dormi-api/internal/dto"
	"dormi-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *model.Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepository) CreateBatch(students []model.Student) error {
	return r.db.CreateInBatches(students, 100).Error
}

func (r *StudentRepository) FindByID(id uuid.UUID) (*model.Student, error) {
	var student model.Student
	err := r.db.First(&student, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) FindByStudentNumber(studentNumber string) (*model.Student, error) {
	var student model.Student
	err := r.db.First(&student, "student_number = ?", studentNumber).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) FindAll(query dto.StudentQuery) ([]model.Student, int64, error) {
	var students []model.Student
	var total int64

	db := r.db.Model(&model.Student{})

	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR student_number ILIKE ?", search, search)
	}
	if query.Grade > 0 {
		db = db.Where("grade = ?", query.Grade)
	}
	if query.Room != "" {
		db = db.Where("room_number = ?", query.Room)
	}

	db.Count(&total)

	offset := (query.Page - 1) * query.Limit
	err := db.Offset(offset).Limit(query.Limit).Order("student_number").Find(&students).Error

	return students, total, err
}

func (r *StudentRepository) Update(student *model.Student) error {
	return r.db.Save(student).Error
}

func (r *StudentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Student{}, "id = ?", id).Error
}

func (r *StudentRepository) ExistsByStudentNumber(studentNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Student{}).Where("student_number = ?", studentNumber).Count(&count).Error
	return count > 0, err
}
