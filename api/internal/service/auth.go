package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"goechotemplate/api/internal/dto"
	"goechotemplate/api/internal/model"
	"goechotemplate/api/internal/repo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	authRepository repo.AuthRepo
}

func NewAuthService(authRepository repo.AuthRepo) AuthService {
	return AuthService{
		authRepository: authRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (dto.LoginResponse, error) {
	person, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("AuthService.Login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("AuthService.Login: %w", err)
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
		return dto.LoginResponse{}, fmt.Errorf("AuthService.Login: %w", err)
	}

	return dto.LoginResponse{
		AccessToken: signedToken,
		AccountID:   person.PersonExternalID,
	}, nil
}

func (s *AuthService) VerifyAuthorisation(
	ctx context.Context,
	params model.VerifyAuthorisationParams,
) (*model.Authorisation, error) {
	if *params.TargetPersonID == *params.AuthenticatedPersonID {
		return &model.Authorisation{}, nil
	}

	return nil, &model.AuthorisationError{}
}
