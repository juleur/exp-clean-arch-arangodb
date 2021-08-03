package session

type Service interface {
	Get(sessionID string) error
	Set(userID string) (string, error)
	Delete(sessionID string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Get(sessionID string) error {
	return s.repository.GetUserSessionID(sessionID)
}

func (s *service) Set(userID string) (string, error) {
	return s.repository.PutUserInSession(userID)
}

func (s *service) Delete(sessionID string) error {
	return s.repository.RemoveUserFromSession(sessionID)
}
