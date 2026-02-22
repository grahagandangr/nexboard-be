package services

import (
	"database/sql"
	"errors"

	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/repositories"
	"github.com/grahagandangr/nexboard-be/utils"
)

type StatusService struct {
	statusRepo *repositories.StatusRepository
}

func NewStatusService(statusRepo *repositories.StatusRepository) *StatusService {
	return &StatusService{statusRepo: statusRepo}
}

// CreateStatus creates new master status
func (s *StatusService) CreateStatus(req *models.StatusRequest) (*models.StatusResponse, error) {
	// Check name uniqueness
	if _, err := s.statusRepo.GetStatusByName(req.Name); err == nil {
		return nil, errors.New("a status with this name already exists")
	}

	pos := 0
	if req.Position != nil {
		pos = *req.Position
	}

	status := &models.Status{
		ExternalID: utils.GenerateUUID(),
		Name:       req.Name,
		Color:      req.Color,
		Position:   pos,
	}

	if err := s.statusRepo.CreateStatus(status); err != nil {
		return nil, err
	}

	return s.mapToResponse(status), nil
}

// GetAllStatuses retrieves all
func (s *StatusService) GetAllStatuses() ([]*models.StatusResponse, error) {
	statuses, err := s.statusRepo.GetAllStatuses()
	if err != nil {
		return nil, err
	}

	var response []*models.StatusResponse
	for _, st := range statuses {
		response = append(response, s.mapToResponse(st))
	}

	return response, nil
}

// GetStatus retrieves detail
func (s *StatusService) GetStatus(externalID string) (*models.StatusResponse, error) {
	status, err := s.statusRepo.GetStatusByExternalID(externalID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("status not found")
		}
		return nil, err
	}

	return s.mapToResponse(status), nil
}

// UpdateStatus modifies status detail
func (s *StatusService) UpdateStatus(externalID string, req *models.StatusRequest) (*models.StatusResponse, error) {
	status, err := s.statusRepo.GetStatusByExternalID(externalID)
	if err != nil {
		return nil, errors.New("status not found")
	}

	// Check if changing name, keep it unique
	if status.Name != req.Name {
		if _, err := s.statusRepo.GetStatusByName(req.Name); err == nil {
			return nil, errors.New("a status with this new name already exists")
		}
	}

	status.Name = req.Name
	status.Color = req.Color
	if req.Position != nil {
		status.Position = *req.Position
	}

	if err := s.statusRepo.UpdateStatus(status); err != nil {
		return nil, err
	}

	return s.mapToResponse(status), nil
}

// DeleteStatus drops status unless referenced
func (s *StatusService) DeleteStatus(externalID string) error {
	status, err := s.statusRepo.GetStatusByExternalID(externalID)
	if err != nil {
		return errors.New("status not found")
	}

	// Check references
	referenced, err := s.statusRepo.CheckIfReferenced(status.ID)
	if err != nil {
		return errors.New("could not verify status references")
	}
	if referenced {
		// As per PRD 7.4 Status Deletion Guard, must return error text mapping to HTTP 409
		return errors.New("conflict: status is in use by one or more tasks and cannot be deleted")
	}

	return s.statusRepo.DeleteStatus(status.ID)
}

// Helper mapper
func (s *StatusService) mapToResponse(st *models.Status) *models.StatusResponse {
	return &models.StatusResponse{
		ExternalID: st.ExternalID,
		Name:       st.Name,
		Color:      st.Color,
		Position:   st.Position,
		CreatedAt:  st.CreatedAt,
		ModifiedAt: st.ModifiedAt,
	}
}
