package aliment

import "exp-clean-arch-arangodb/pkg/entities"

type Service interface {
	GetCategories() ([]string, error)
	Create(userID string, aliment entities.AlimentCollection) error
	Delete(userID string, alimentID string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetCategories() ([]string, error) {
	return s.repository.FetchCategories()
}

func (s *service) Create(userID string, aliment entities.AlimentCollection) error {
	return s.repository.InsertAliment(userID, aliment)
}

func (s *service) Delete(userID string, alimentID string) error {
	return s.repository.DeleteAliment(userID, alimentID)
}
