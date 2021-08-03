package city

import "exp-clean-arch-arangodb/pkg/entities"

type Service interface {
	FindCity(cityName string) ([]string, error)
	IsCityExist(cityName string) (*entities.CityCollection, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FindCity(cityName string) ([]string, error) {
	return s.repository.GetCity(cityName)
}

func (s *service) IsCityExist(cityName string) (*entities.CityCollection, error) {
	return s.repository.CheckIfCityNameExists(cityName)
}
