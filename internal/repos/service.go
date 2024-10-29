package repos

import "justlend/internal/tron"

type Service struct {
	tron *tron.Endpoint
}

func NewService(endpoint *tron.Endpoint) *Service {
	return &Service{
		tron: endpoint,
	}
}
