package repository

import (
	"context"
	"database/sql"
	db "goechotemplate/api/db/model"
	"goechotemplate/api/internal/account/model"
	"time"
)

type PersonRepository interface {
	GetByExternalID(ctx context.Context, externalID string) (*model.Person, error)
	Create(ctx context.Context, person *model.Person) (*model.Person, error)
}

type personRepository struct {
	queries *db.Queries
}

func NewPersonRepository(dbtx db.DBTX) PersonRepository {
	return &personRepository{
		queries: db.New(dbtx),
	}
}

func (r *personRepository) GetByExternalID(ctx context.Context, externalID string) (*model.Person, error) {
	person, err := r.queries.GetPerson(ctx, externalID)
	if err != nil {
		return &model.Person{}, err
	}
	return &model.Person{
		ID:         person.ID,
		ExternalID: person.ExternalID,
		Email:      person.Email.String,
		Password:   person.Password,
		CreatedAt:  person.CreatedAt.Time,
		UpdatedAt:  person.UpdatedAt.Time,
	}, nil
}

func (r *personRepository) Create(ctx context.Context, person *model.Person) (*model.Person, error) {
	createdPerson, err := r.queries.CreatePerson(ctx, db.CreatePersonParams{
		ExternalID: person.ExternalID,
		Email:      sql.NullString{String: person.Email, Valid: true},
		Password:   nil,
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return &model.Person{}, err
	}

	return &model.Person{
		ID:         createdPerson.ID,
		ExternalID: createdPerson.ExternalID,
		Email:      createdPerson.Email.String,
		Password:   createdPerson.Password,
		CreatedAt:  createdPerson.CreatedAt.Time,
		UpdatedAt:  createdPerson.UpdatedAt.Time,
	}, nil
}
