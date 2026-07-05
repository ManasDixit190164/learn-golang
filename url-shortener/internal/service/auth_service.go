package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manasdixit/url-shortener/internal/domain"
	"github.com/manasdixit/url-shortener/internal/repository"
	"github.com/manasdixit/url-shortener/internal/utils"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	userRepo           repository.UserRepository
	refreshTokenRepo   repository.RefreshTokenRepository
	jwtManager         *utils.JWTManager
	refreshTokenExpiry time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *utils.JWTManager,
	refreshTokenExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		jwtManager:         jwtManager,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (s *AuthService) Signup(ctx context.Context, req domain.SignupRequest) (domain.AuthResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Name == "" || !utils.IsValidEmail(req.Email) || len(req.Password) < 8 {
		return domain.AuthResponse{}, ErrInvalidInput
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	user := domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.AuthResponse{}, ErrEmailAlreadyExists
		}
		return domain.AuthResponse{}, err
	}

	return s.createAuthResponse(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (domain.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.AuthResponse{}, ErrInvalidCredentials
		}
		return domain.AuthResponse{}, err
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return domain.AuthResponse{}, ErrInvalidCredentials
	}

	return s.createAuthResponse(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (domain.AuthResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return domain.AuthResponse{}, ErrInvalidToken
	}

	tokenHash := utils.HashToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.GetValidByHash(ctx, tokenHash)
	if err != nil {
		return domain.AuthResponse{}, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return domain.AuthResponse{}, ErrInvalidToken
	}

	_ = s.refreshTokenRepo.RevokeByHash(ctx, tokenHash)

	return s.createAuthResponse(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return ErrInvalidToken
	}

	tokenHash := utils.HashToken(refreshToken)
	if err := s.refreshTokenRepo.RevokeByHash(ctx, tokenHash); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidToken
		}
		return err
	}

	return nil
}

func (s *AuthService) createAuthResponse(ctx context.Context, user domain.User) (domain.AuthResponse, error) {
	accessToken, err := s.jwtManager.Generate(user.ID, user.Email)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	refreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	_, err = s.refreshTokenRepo.Create(ctx, domain.RefreshToken{
		ID:        uuid.Nil,
		UserID:    user.ID,
		TokenHash: utils.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.refreshTokenExpiry),
	})
	if err != nil {
		return domain.AuthResponse{}, err
	}

	return domain.AuthResponse{
		User:         domain.NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
