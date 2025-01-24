package services

import (
	"errors"
	"server/internal/models"
	"server/internal/repositories"
)

type FacilityService struct {
	repo repositories.FacilityRepository
}

// NewFacilityService initializes a new FacilityService.
func NewFacilityService(repo repositories.FacilityRepository) *FacilityService {
	return &FacilityService{repo: repo}
}

func (s *FacilityService) GetAllFacilities() (*[]models.Facility, error) {
	facilites, err := s.repo.FindMany((map[string]interface{}{}))
	if err != nil {
		return nil, err
	}

	return &facilites, nil
}

// GetFacilityByID fetches a facility by its ID.
func (s *FacilityService) GetFacilityByID(id int64) (*models.Facility, error) {
	facility, err := s.repo.Find(id)
	if err != nil {
		return nil, errors.New("facility not found")
	}
	return facility, nil
}
