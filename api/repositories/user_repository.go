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

func (repo *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := repo.db.QueryRow("SELECT id, password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) Create(user models.User) error {
	_, err := repo.db.Exec("INSERT INTO users (id, password) VALUES (?, ?)", user.ID, user.Password)
	return err
}
