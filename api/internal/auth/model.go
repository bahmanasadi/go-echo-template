package auth

import echojwt "github.com/labstack/echo-jwt/v4"

var DefaultJWTConfig = echojwt.Config{
	ContextKey: "user",
	SigningKey: []byte("secret"),
}

type AuthPerson struct {
	PersonExternalID string
	Password         []byte
}

type VerifyAuthorisationParams struct {
	AuthenticatedPersonID *string
	TargetPersonID        *string
	OrganisationID        *string
}

type Authorisation struct {
}

type AuthorisationError struct {
	msg string
}

func (e *AuthorisationError) Error() string { return e.msg }
