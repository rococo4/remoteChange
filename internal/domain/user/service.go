package user

import (
	"fmt"
	"remoteChange/internal/domain/jwt"
	"remoteChange/internal/infrastructure"
	"remoteChange/internal/model"
)

type Service struct {
	repo userRepo
}

func NewService(repo userRepo) Service {
	return Service{repo: repo}
}

func (s *Service) Register(user model.UserDTORegister) (string, error) {
	entity := infrastructure.MapUserDtoRegisterToUserEntity(user)
	token, err := jwt.GenerateJWT(entity.Username, entity.Role)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}
	entity.Password, err = infrastructure.HashPassword(entity.Password)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	err = s.repo.SaveUser(entity)
	if err != nil {
		return "", fmt.Errorf("error saving user: %w", err)
	}
	return token, nil
}

func (s *Service) Login(user model.UserDTOLogin) (string, error) {
	entity, err := s.repo.GetUserByUsername(user.Username)
	if err != nil {
		return "", fmt.Errorf("error generating user: %w", err)
	}
	if !infrastructure.CheckPassword(user.Password, entity.Password) {
		return "", fmt.Errorf("passwords do not match")
	}
	token, err := jwt.GenerateJWT(entity.Username, entity.Role)
	return token, nil
}
