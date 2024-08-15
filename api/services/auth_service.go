package services

import (
	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) ValidateUser(user models.User) bool {
	storedUser, err := s.repo.GetByID(user.ID)
	if err != nil {
		return false
	}
	return storedUser.Password == user.Password
}
