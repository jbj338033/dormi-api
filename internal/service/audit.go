package service

import (
	"encoding/json"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"
	"dormi-api/internal/repository"

	"github.com/google/uuid"
)

type AuditService struct {
	auditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

func (s *AuditService) Log(userID uuid.UUID, action model.AuditAction, entityType string, entityID *uuid.UUID, details interface{}, ipAddress string) error {
	var detailsJSON []byte
	if details != nil {
		var err error
		detailsJSON, err = json.Marshal(details)
		if err != nil {
			detailsJSON = nil
		}
	}

	log := &model.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    detailsJSON,
		IPAddress:  ipAddress,
	}

	return s.auditRepo.Create(log)
}

func (s *AuditService) GetAll(query dto.AuditQuery) ([]model.AuditLog, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	return s.auditRepo.FindAll(query)
}
