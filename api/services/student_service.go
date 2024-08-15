package services

import (
	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
)

type StudentService struct {
	repo *repositories.StudentRepository
}

func NewStudentService(repo *repositories.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	return s.repo.GetAll()
}

func (s *StudentService) GetStudentByID(id string) (*models.Student, error) {
	return s.repo.GetByID(id)
}

func (s *StudentService) CreateStudent(student models.Student) error {
	return s.repo.Create(student)
}

func (s *StudentService) UpdateStudent(id string, student models.Student) error {
	return s.repo.Update(id, student)
}

func (s *StudentService) DeleteStudent(id string) error {
	return s.repo.Delete(id)
}
