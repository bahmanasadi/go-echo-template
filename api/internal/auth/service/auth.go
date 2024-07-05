package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"goechotemplate/api/internal/auth/dto"
	"goechotemplate/api/internal/auth/model"
	"goechotemplate/api/internal/auth/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (dto.LoginResponse, error)
	VerifyAuthorisation(ctx context.Context, params model.VerifyAuthorisationParams) (*model.Authorisation, error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &authService{
		authRepository: authRepository,
	}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (dto.LoginResponse, error) {
	person, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, err
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
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		AccessToken: signedToken,
		AccountID:   person.PersonExternalID,
	}, nil
}

func (s *authService) VerifyAuthorisation(
	ctx context.Context,
	params model.VerifyAuthorisationParams,
) (*model.Authorisation, error) {
	if *params.TargetPersonID == *params.AuthenticatedPersonID {
		return &model.Authorisation{}, nil
	}

	return nil, &model.AuthorisationError{}
}
