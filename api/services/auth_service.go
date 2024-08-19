package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) ValidateUser(ctx context.Context, email, password string) (bool, error) {
	// this should return true if the credentials are valid
	return s.repo.ValidateCredentials(ctx, email, password)
}

func (s *AuthService) GetUserIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// log in the log file
		return uuid.UUID{}, errors.New("user not found")
	}
	return user.ID, nil
}

func (s *AuthService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	exits, err := s.repo.UserExistsByEmail(ctx, user.Email)
	if exits || err != nil {
		return nil, errors.New("user already exists")
	}
	user.ID = uuid.New()

	return s.repo.CreateUser(ctx, user)
}
