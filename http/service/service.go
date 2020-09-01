package service

type Service struct {
	manager QueryManager
}

func NewService(manager QueryManager) *Service {
	return &Service{manager: manager}
}
