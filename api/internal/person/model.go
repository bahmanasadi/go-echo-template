package person

import (
	"time"
)

type Person struct {
	ID         int64     `json:"id"`
	ExternalID string    `json:"externalId"`
	Email      string    `json:"email" validate:"required,email"`
	Password   []byte    `json:"password"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
