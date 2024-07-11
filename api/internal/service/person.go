package service

import (
	"context"
	"fmt"
	"goechotemplate/api/internal/model"
	"goechotemplate/api/internal/repo"
)

type PersonService struct {
	personRepo repo.PersonRepo
}

func NewPersonService(personRepo repo.PersonRepo) PersonService {
	return PersonService{
		personRepo: personRepo,
	}
}

func (s *PersonService) GetByExternalID(ctx context.Context, externalID string) (*model.Person, error) {
	p, err := s.personRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, fmt.Errorf("PersonService.GetByExternalID: %w", err)
	}
	return p, nil
}

func (s *PersonService) Create(ctx context.Context, person *model.Person) (*model.Person, error) {
	createdPerson, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("PersonService.Create: %w", err)
	}
	return createdPerson, nil
}
