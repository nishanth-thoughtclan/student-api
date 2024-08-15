package repositories

import (
	"database/sql"

	"github.com/nishanth-thoughtclan/student-api/api/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (repo *StudentRepository) GetAll() ([]models.Student, error) {
	rows, err := repo.db.Query("SELECT id, name, age, created_by, created_on, updated_by, updated_on FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.CreatedBy, &student.CreatedOn, &student.UpdatedBy, &student.UpdatedOn); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (repo *StudentRepository) GetByID(id string) (*models.Student, error) {
	var student models.Student
	err := repo.db.QueryRow("SELECT id, name, age, created_by, created_on, updated_by, updated_on FROM students WHERE id = ?", id).Scan(&student.ID, &student.Name, &student.Age, &student.CreatedBy, &student.CreatedOn, &student.UpdatedBy, &student.UpdatedOn)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (repo *StudentRepository) Create(student models.Student) error {
	_, err := repo.db.Exec("INSERT INTO students (id, name, age, created_by, created_on, updated_by, updated_on) VALUES (?, ?, ?, ?, ?, ?, ?)", student.ID, student.Name, student.Age, student.CreatedBy, student.CreatedOn, student.UpdatedBy, student.UpdatedOn)
	return err
}

func (repo *StudentRepository) Update(id string, student models.Student) error {
	_, err := repo.db.Exec("UPDATE students SET name = ?, age = ?, updated_by = ?, updated_on = ? WHERE id = ?", student.Name, student.Age, student.UpdatedBy, student.UpdatedOn, id)
	return err
}

func (repo *StudentRepository) Delete(id string) error {
	_, err := repo.db.Exec("DELETE FROM students WHERE id = ?", id)
	return err
}
