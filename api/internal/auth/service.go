package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	authRepository Repository
}

func NewAuthService(authRepository Repository) Service {
	return Service{
		authRepository: authRepository,
	}
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (LoginResponse, error) {
	person, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Service.Login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(req.Password))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Service.Login: %w", err)
	}

	claims := jwt.RegisteredClaims{
		Subject:   person.PersonExternalID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Service.Login: %w", err)
	}

	return LoginResponse{
		AccessToken: signedToken,
		AccountID:   person.PersonExternalID,
	}, nil
}

func (s *Service) VerifyAuthorisation(
	ctx context.Context,
	params VerifyAuthorisationParams,
) (*Authorisation, error) {
	if *params.TargetPersonID == *params.AuthenticatedPersonID {
		return &Authorisation{}, nil
	}

	return nil, &AuthorisationError{}
}
