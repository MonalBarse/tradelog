package service

import (
	"context"
	"errors"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/MonalBarse/tradelog/internal/repository"
	"github.com/MonalBarse/tradelog/pkg/utils"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

// @desc: reg new user
// @flow: check existing -> hash pwd -> create user
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
		Role:     "user", // default-> user
	}

	return s.repo.Create(ctx, user)
}

// @desc: login user
// @flow: verify email -> check password -> generate tokens
func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate both tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
