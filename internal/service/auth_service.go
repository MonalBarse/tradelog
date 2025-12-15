package service

import (
	"context"
	"errors"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/MonalBarse/tradelog/internal/repository"
	"github.com/MonalBarse/tradelog/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*domain.User, string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	PromoteToAdmin(ctx context.Context, userID uint, secret string) error
}

type authService struct {
	repo          repository.UserRepository
	jwtSecret     string
	refreshSecret string
	adminSecret   string
}

func NewAuthService(repo repository.UserRepository, jwtSecret, refreshSecret, adminSecret string) AuthService {
	return &authService{
		repo:          repo,
		jwtSecret:     jwtSecret,
		refreshSecret: refreshSecret,
		adminSecret:   adminSecret,
	}
}

// @desc: register new user
// @flow: check existing user -> hash password -> create user record
func (s *authService) Register(ctx context.Context, email, password string) error {
	existingUser, _ := s.repo.FindByEmail(ctx, email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
	}

	return s.repo.Create(ctx, user)
}

// @desc: login user
// @flow: find user by email -> verify password -> generate tokens
func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Role, s.jwtSecret, s.refreshSecret)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// @desc: refresh tokens
// @flow: validate refresh token -> generate new tokens
func (s *authService) Refresh(ctx context.Context, refreshTokenString string) (string, string, error) {
	token, err := utils.ValidateRefreshToken(refreshTokenString, s.refreshSecret)
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID := uint(claims["sub"].(float64))

	// In a real app, I might query DB here to check user status using userID
	// user, err := s.repo.FindByID(ctx, userID) ...

	return utils.GenerateTokens(userID, "user", s.jwtSecret, s.refreshSecret)
}

func (s *authService) PromoteToAdmin(ctx context.Context, userID uint, secret string) error {
	if secret != s.adminSecret {
		return errors.New("invalid admin secret")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Role = "admin"
	return s.repo.Update(ctx, user)
}
