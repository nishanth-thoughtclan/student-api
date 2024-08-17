package repositories

import (
	"database/sql"

	"github.com/nishanth-thoughtclan/student-api/api/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// func (repo *UserRepository) Create(user models.User) error {
// 	_, err := repo.db.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", user.ID, user.Email, user.Password)
// 	return err
// }

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ValidateCredentials(email, password string) (bool, error) {
	// Logic to validate user credentials
	var storedPassword string
	err := r.db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil || storedPassword != password {
		return false, err
	}
	return true, nil
}
