package person

import (
	"context"
	"database/sql"
	"fmt"
	db "goechotemplate/api/db/model"
	"time"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return Repository{queries: queries}
}

func (r *Repository) GetByExternalID(ctx context.Context, externalID string) (*Person, error) {
	person, err := r.queries.GetPerson(ctx, externalID)
	if err != nil {
		return &Person{}, fmt.Errorf("Repository.GetByExternalID: %w", err)
	}
	return &Person{
		ID:         person.ID,
		ExternalID: person.ExternalID,
		Email:      person.Email.String,
		Password:   person.Password,
		CreatedAt:  person.CreatedAt.Time,
		UpdatedAt:  person.UpdatedAt.Time,
	}, nil
}

func (r *Repository) Create(ctx context.Context, person *Person) (*Person, error) {
	createdPerson, err := r.queries.CreatePerson(ctx, db.CreatePersonParams{
		ExternalID: person.ExternalID,
		Email:      sql.NullString{String: person.Email, Valid: true},
		Password:   nil,
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return &Person{}, fmt.Errorf("Repository.Create: %w", err)
	}

	return &Person{
		ID:         createdPerson.ID,
		ExternalID: createdPerson.ExternalID,
		Email:      createdPerson.Email.String,
		Password:   createdPerson.Password,
		CreatedAt:  createdPerson.CreatedAt.Time,
		UpdatedAt:  createdPerson.UpdatedAt.Time,
	}, nil
}
