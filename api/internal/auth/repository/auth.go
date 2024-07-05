package repository

import (
	"context"
	"database/sql"
	"fmt"
	db "goechotemplate/api/db/model"
	"goechotemplate/api/internal/auth/model"
)

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (model.AuthPerson, error)
}

type authRepository struct {
	queries *db.Queries
}

func NewAuthRepository(queries *db.Queries) AuthRepository {
	return &authRepository{queries: queries}
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (model.AuthPerson, error) {
	person, err := r.queries.GetPersonByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		return model.AuthPerson{}, fmt.Errorf("GetByEmail: %w", err)
	}

	return model.AuthPerson{
		PersonExternalID: person.ExternalID,
		Password:         person.Password,
	}, nil
}
