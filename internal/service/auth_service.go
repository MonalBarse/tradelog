package service

import (
	"errors"

	"github.com/MonalBarse/tradelog/internal/domain"
	"github.com/MonalBarse/tradelog/internal/repository"
	"github.com/MonalBarse/tradelog/pkg/utils"
)

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

// @desc: reg new user
// @flow : check if user exists -> if not, hash pswrd -> create user
func (s *authService) Register(email, password string) error {

	existingUser, _ := s.repo.FindByEmail(email)
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
		Role:     "user", //default-> user
	}

	return s.repo.Create(user)
}

// @desc: login user
// @flow: find user by email -> check pswrd -> generate jwt
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	// if no err -> user found move ahead
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}