package person

import (
	"context"
	"fmt"
)

type Service struct {
	personRepo Repository
}

func NewService(personRepo Repository) Service {
	return Service{
		personRepo: personRepo,
	}
}

func (s *Service) GetByExternalID(ctx context.Context, externalID string) (*Person, error) {
	p, err := s.personRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, fmt.Errorf("Service.GetByExternalID: %w", err)
	}
	return p, nil
}

func (s *Service) Create(ctx context.Context, person *Person) (*Person, error) {
	createdPerson, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("Service.Create: %w", err)
	}
	return createdPerson, nil
}
