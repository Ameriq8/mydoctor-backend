package services

import (
	"errors"
	"server/internal/models"
	"server/internal/repositories"
)

type CityService struct {
	repo repositories.CitiesRepository
}

func NewCityService(repo repositories.CitiesRepository) *CityService {
	return &CityService{repo: repo}
}

func (s *CityService) GetCityByID(id int64) (*models.City, error) {
	city, err := s.repo.Find(id)
	if err != nil {
		return nil, errors.New("city not found")
	}
	return city, nil
}
