// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package model

import (
	"database/sql"
)

type Person struct {
	ID         int64
	ExternalID string
	Email      sql.NullString
	Password   []byte
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}
