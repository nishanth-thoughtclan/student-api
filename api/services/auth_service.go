package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) ValidateUser(email, password string) (bool, error) {
	// Logic to validate user credentials
	// This should return true if the credentials are valid
	return s.repo.ValidateCredentials(email, password)
}

func (s *AuthService) GetUserIDByEmail(email string) (uuid.UUID, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		// log in the log file
		return uuid.UUID{}, errors.New("user not found")
	}
	return user.ID, nil
}
