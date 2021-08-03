package activity

type Service interface {
	LastUserActivity(sessionID string)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) LastUserActivity(sessionID string) {
	go s.repository.UpdateLastSeenUser(sessionID)
}
