package service // package declaration for the module

import ( // start import block
	"context" // import package
	"errors" // import package
	"strings" // import package
	"time" // import package

	"github.com/google/uuid" // import package
	"github.com/manasdixit/url-shortener/internal/domain" // import package
	"github.com/manasdixit/url-shortener/internal/repository" // import package
	"github.com/manasdixit/url-shortener/internal/utils" // import package
) // end import block or block scope

var ( // declare variable
	ErrInvalidInput       = errors.New("invalid input") // assign value
	ErrInvalidCredentials = errors.New("invalid email or password") // assign value
	ErrEmailAlreadyExists = errors.New("email already exists") // assign value
	ErrInvalidToken       = errors.New("invalid token") // assign value
) // end import block or block scope

type AuthService struct { // declare struct type
	userRepo           repository.UserRepository // execute statement
	refreshTokenRepo   repository.RefreshTokenRepository // execute statement
	jwtManager         *utils.JWTManager // execute statement
	refreshTokenExpiry time.Duration // execute statement
} // end block

func NewAuthService( // declare function
	userRepo repository.UserRepository, // execute statement
	refreshTokenRepo repository.RefreshTokenRepository, // execute statement
	jwtManager *utils.JWTManager, // execute statement
	refreshTokenExpiry time.Duration, // execute statement
) *AuthService { // execute statement
	return &AuthService{ // return statement
		userRepo:           userRepo, // execute statement
		refreshTokenRepo:   refreshTokenRepo, // execute statement
		jwtManager:         jwtManager, // execute statement
		refreshTokenExpiry: refreshTokenExpiry, // execute statement
	} // end block
} // end block

func (s *AuthService) Signup(ctx context.Context, req domain.SignupRequest) (domain.AuthResponse, error) { // declare method
	req.Name = strings.TrimSpace(req.Name) // trim whitespace
	req.Email = strings.ToLower(strings.TrimSpace(req.Email)) // trim whitespace

	if req.Name == "" || !utils.IsValidEmail(req.Email) || len(req.Password) < 8 { // if condition
		return domain.AuthResponse{}, ErrInvalidInput // return statement
	} // end block

	passwordHash, err := utils.HashPassword(req.Password) // hash user password securely
	if err != nil { // if condition
		return domain.AuthResponse{}, err // return statement
	} // end block

	user := domain.User{ // declare and initialize variable
		Name:         req.Name, // execute statement
		Email:        req.Email, // execute statement
		PasswordHash: passwordHash, // execute statement
	} // end block

	user, err = s.userRepo.Create(ctx, user) // assign value
	if err != nil { // if condition
		if errors.Is(err, repository.ErrConflict) { // handle duplicate database entry
			return domain.AuthResponse{}, ErrEmailAlreadyExists // return statement
		} // end block
		return domain.AuthResponse{}, err // return statement
	} // end block

	return s.createAuthResponse(ctx, user) // return statement
} // end block

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (domain.AuthResponse, error) { // declare method
	email := strings.ToLower(strings.TrimSpace(req.Email)) // trim whitespace

	user, err := s.userRepo.GetByEmail(ctx, email) // declare and initialize variable
	if err != nil { // if condition
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return domain.AuthResponse{}, ErrInvalidCredentials // return statement
		} // end block
		return domain.AuthResponse{}, err // return statement
	} // end block

	if !utils.CheckPassword(req.Password, user.PasswordHash) { // compare password with hash
		return domain.AuthResponse{}, ErrInvalidCredentials // return statement
	} // end block

	return s.createAuthResponse(ctx, user) // return statement
} // end block

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (domain.AuthResponse, error) { // declare method
	if strings.TrimSpace(refreshToken) == "" { // trim whitespace
		return domain.AuthResponse{}, ErrInvalidToken // return statement
	} // end block

	tokenHash := utils.HashToken(refreshToken) // hash refresh token for storage
	storedToken, err := s.refreshTokenRepo.GetValidByHash(ctx, tokenHash) // declare and initialize variable
	if err != nil { // if condition
		return domain.AuthResponse{}, ErrInvalidToken // return statement
	} // end block

	user, err := s.userRepo.GetByID(ctx, storedToken.UserID) // declare and initialize variable
	if err != nil { // if condition
		return domain.AuthResponse{}, ErrInvalidToken // return statement
	} // end block

	_ = s.refreshTokenRepo.RevokeByHash(ctx, tokenHash) // assign value

	return s.createAuthResponse(ctx, user) // return statement
} // end block

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error { // declare method
	if strings.TrimSpace(refreshToken) == "" { // trim whitespace
		return ErrInvalidToken // return statement
	} // end block

	tokenHash := utils.HashToken(refreshToken) // hash refresh token for storage
	if err := s.refreshTokenRepo.RevokeByHash(ctx, tokenHash); err != nil { // if condition
		if errors.Is(err, repository.ErrNotFound) { // handle missing database record
			return ErrInvalidToken // return statement
		} // end block
		return err // return statement
	} // end block

	return nil // return statement
} // end block

func (s *AuthService) createAuthResponse(ctx context.Context, user domain.User) (domain.AuthResponse, error) { // declare method
	accessToken, err := s.jwtManager.Generate(user.ID, user.Email) // declare and initialize variable
	if err != nil { // if condition
		return domain.AuthResponse{}, err // return statement
	} // end block

	refreshToken, err := utils.GenerateSecureToken(32) // create a random secure token
	if err != nil { // if condition
		return domain.AuthResponse{}, err // return statement
	} // end block

	_, err = s.refreshTokenRepo.Create(ctx, domain.RefreshToken{ // assign value
		ID:        uuid.Nil, // execute statement
		UserID:    user.ID, // execute statement
		TokenHash: utils.HashToken(refreshToken), // hash refresh token for storage
		ExpiresAt: time.Now().Add(s.refreshTokenExpiry), // execute statement
	}) // close block
	if err != nil { // if condition
		return domain.AuthResponse{}, err // return statement
	} // end block

	return domain.AuthResponse{ // return statement
		User:         domain.NewUserResponse(user), // execute statement
		AccessToken:  accessToken, // execute statement
		RefreshToken: refreshToken, // execute statement
	}, nil // close block
} // end block
