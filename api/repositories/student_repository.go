package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/nishanth-thoughtclan/student-api/api/models"
)

var dateTimeFormat string = "2006-01-02 15:04:05"

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
		var createdOn, updatedOn string

		// Scan the result into the student object
		if err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.CreatedBy, &createdOn, &student.UpdatedBy, &updatedOn); err != nil {
			return nil, err
		}

		// Parsing the time values from the database using a format
		student.CreatedOn, err = time.Parse(dateTimeFormat, createdOn)
		if err != nil {
			return nil, err
		}
		student.UpdatedOn, err = time.Parse(dateTimeFormat, updatedOn)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}
	return students, nil
}

func (repo *StudentRepository) GetByID(ctx context.Context, id string) (*models.Student, error) {
	var student models.Student
	var createdOn, updatedOn string

	// Query and scan into the appropriate variables
	err := repo.db.QueryRowContext(ctx, "SELECT id, name, age, created_by, created_on, updated_by, updated_on FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.Name, &student.Age, &student.CreatedBy, &createdOn, &student.UpdatedBy, &updatedOn)
	if err != nil {
		return nil, err
	}

	// Parse the string dates into time.Time
	student.CreatedOn, err = time.Parse(dateTimeFormat, createdOn)
	if err != nil {
		return nil, err
	}
	student.UpdatedOn, err = time.Parse(dateTimeFormat, updatedOn)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (repo *StudentRepository) Create(ctx context.Context, student models.Student) (*models.Student, error) {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age, created_by, created_on, updated_by, updated_on) VALUES (?, ?, ?, ?, ?, ?, ?)", student.ID, student.Name, student.Age, student.CreatedBy, student.CreatedOn, student.UpdatedBy, student.UpdatedOn)
	if err != nil {
		return &models.Student{}, err
	}

	return &student, err
}

func (repo *StudentRepository) Update(ctx context.Context, id string, student models.Student) (*models.Student, error) {
	_, err := repo.db.ExecContext(ctx, "UPDATE students SET name = ?, age = ?, updated_by = ?, updated_on = ? WHERE id = ?", student.Name, student.Age, student.UpdatedBy, student.UpdatedOn, id)
	return &student, err
}

func (repo *StudentRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM students WHERE id = ?", id)
	return err
}
