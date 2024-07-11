package auth

import (
	"context"
	"database/sql"
	"fmt"
	db "goechotemplate/api/db/model"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return Repository{queries: queries}
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (AuthPerson, error) {
	person, err := r.queries.GetPersonByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		return AuthPerson{}, fmt.Errorf("GetByEmail: %w", err)
	}

	return AuthPerson{
		PersonExternalID: person.ExternalID,
		Password:         person.Password,
	}, nil
}
