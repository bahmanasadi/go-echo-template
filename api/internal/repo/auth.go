package repo

import (
	"context"
	"database/sql"
	"fmt"
	db "goechotemplate/api/db/model"
	"goechotemplate/api/internal/model"
)

type AuthRepo struct {
	queries *db.Queries
}

func NewAuthRepo(queries *db.Queries) AuthRepo {
	return AuthRepo{queries: queries}
}

func (r *AuthRepo) GetByEmail(ctx context.Context, email string) (model.AuthPerson, error) {
	person, err := r.queries.GetPersonByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		return model.AuthPerson{}, fmt.Errorf("GetByEmail: %w", err)
	}

	return model.AuthPerson{
		PersonExternalID: person.ExternalID,
		Password:         person.Password,
	}, nil
}
