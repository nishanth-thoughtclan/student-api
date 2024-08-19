package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nishanth-thoughtclan/student-api/api/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ValidateCredentials(ctx context.Context, email, password string) (bool, error) {
	var user models.User
	var storedPassword string
	err := r.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	// check if the provided password matches the hashed password
	if err != nil || !user.CheckPassword(password, storedPassword) {
		return false, errors.New("invalid email/password")
	}
	return true, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	// Hash the user's password before storing it
	if err := user.HashPassword(); err != nil {
		return nil, err
	}
	// Then save the user to the database
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES (?, ?, ?)", user.ID, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? LIMIT 1)"
	err := repo.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
