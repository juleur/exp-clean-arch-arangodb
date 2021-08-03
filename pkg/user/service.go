package user

import (
	"exp-clean-arch-arangodb/pkg/entities"
)

type Service interface {
	Create(newUser entities.NewUserBody) (string, error)
	FindByEmail(email string) (*entities.UserDoc, error)
	LastSeenAt(userID string) error
	AddCity(userID string, cityID string) error
	DeleteCity(userID string) error
	DeleteTotallyAccount(userID string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(newUser entities.NewUserBody) (string, error) {
	return s.repository.InsertUser(newUser)
}

func (s *service) FindByEmail(email string) (*entities.UserDoc, error) {
	return s.repository.FetchUserByEmail(email)
}

func (s *service) LastSeenAt(userID string) error {
	return s.repository.UpdateLastSeenUser(userID)
}

func (s *service) AddCity(userID string, cityID string) error {
	return s.repository.InsertCityUser(userID, cityID)
}

func (s *service) DeleteCity(userID string) error {
	return s.repository.DeleteCityUser(userID)
}

func (s *service) DeleteTotallyAccount(userID string) error {
	return s.repository.DeleteTotallyUser(userID)
}
