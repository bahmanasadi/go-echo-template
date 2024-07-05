package service

import (
	"context"
	"fmt"
	"goechotemplate/api/internal/account/model"
	"goechotemplate/api/internal/account/repository"
)

type PersonService interface {
	GetPersonByExternalID(ctx context.Context, externalID string) (*model.Person, error)
	CreatePerson(ctx context.Context, person *model.Person) (*model.Person, error)
}

type personService struct {
	personRepo repository.PersonRepository
}

func NewPersonService(personRepo repository.PersonRepository) PersonService {
	return &personService{
		personRepo: personRepo,
	}
}

func (s *personService) GetPersonByExternalID(ctx context.Context, externalID string) (*model.Person, error) {
	p, err := s.personRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, fmt.Errorf("GetPersonByExternalID: %w", err)
	}
	return p, nil
}

func (s *personService) CreatePerson(ctx context.Context, person *model.Person) (*model.Person, error) {
	createdPerson, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("CreatePerson: %w", err)
	}
	return createdPerson, nil
}
