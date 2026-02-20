package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"ticketing/backend/internal/config"
	"ticketing/backend/internal/domain"
	"ticketing/backend/internal/dto"
	"ticketing/backend/internal/repository"
)

// UserService handles registration, login, and JWT issuance.
type UserService struct {
	userRepo repository.UserRepository
	jwt      config.JWTConfig
}

func NewUserService(repo repository.UserRepository, jwtCfg config.JWTConfig) *UserService {
	return &UserService{userRepo: repo, jwt: jwtCfg}
}

// Register creates a new user account with a bcrypt-hashed password.
//
// bcrypt design decisions:
//   - Cost factor 12: ~250ms on modern hardware — enough to slow brute-force attacks
//     without impacting interactive login UX.
//   - bcrypt incorporates a random salt automatically; we never store the salt separately.
//   - The hash is stored as a self-contained string: "$2a$12$<22-char-salt><31-char-hash>"
func (s *UserService) Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.LoginResponse, error) {
	// Validate role
	role := domain.UserRole(req.Role)
	if role != domain.RoleAttendee && role != domain.RoleOrganizer && role != domain.RoleAdmin {
		return nil, fmt.Errorf("invalid role: %s", req.Role)
	}

	// Hash password — bcrypt cost 12 is the production-safe minimum
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("UserService.Register: hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err // propagates ErrDuplicateEmail
	}

	return s.buildLoginResponse(user)
}

// Login verifies credentials and returns a signed JWT on success.
func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Return the same error for wrong email and wrong password
			// to prevent email enumeration attacks.
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("UserService.Login: lookup: %w", err)
	}

	// bcrypt.CompareHashAndPassword is timing-safe — it always takes ~250ms
	// regardless of whether the password matches, preventing timing attacks.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, repository.ErrNotFound // Same error — no enumeration
	}

	return s.buildLoginResponse(user)
}

// GetProfile fetches a user's public profile.
func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// buildLoginResponse constructs a JWT and response payload for a given user.
func (s *UserService) buildLoginResponse(user *domain.User) (*dto.LoginResponse, error) {
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserProfile{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
		},
	}, nil
}

func (s *UserService) generateToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  string(user.Role),
		"iat":   now.Unix(),
		"exp":   now.Add(time.Duration(s.jwt.ExpiryMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwt.Secret))
	if err != nil {
		return "", fmt.Errorf("generateToken: sign: %w", err)
	}
	return signed, nil
}
