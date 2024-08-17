package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
)

type StudentService struct {
	repo *repositories.StudentRepository
}

func NewStudentService(repo *repositories.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	return s.repo.GetAll()
}

func (s *StudentService) GetStudentByID(ctx context.Context, id string) (*models.Student, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *StudentService) CreateStudent(ctx context.Context, student models.Student) (*models.Student, error) {
	currentUserId := ctx.Value("userID").(string)
	student.ID = uuid.New()
	userID, err := uuid.Parse(currentUserId)
	if err != nil {
		return &models.Student{}, errors.New("invalid user ID format")
	}
	student.CreatedBy = userID
	student.UpdatedBy = userID
	student.CreatedOn = time.Now()
	student.UpdatedOn = time.Now()
	return s.repo.Create(ctx, student)
}

func (s *StudentService) UpdateStudent(ctx context.Context, id string, updatedStudent models.Student) (*models.Student, error) {
	currentUserId := ctx.Value("userID").(string)
	existingStudent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &models.Student{}, errors.New("student not found")
	}

	if existingStudent.CreatedBy.String() != currentUserId {
		return &models.Student{}, errors.New("you are not allowed to update this student")
	}

	parsedUserID, err := uuid.Parse(currentUserId)
	if err != nil {
		return &models.Student{}, errors.New("invalid user ID format")
	}
	// Preserving the existing ID and CreatedOn values
	updatedStudent.ID = existingStudent.ID
	updatedStudent.CreatedBy = existingStudent.CreatedBy
	updatedStudent.CreatedOn = existingStudent.CreatedOn

	updatedStudent.UpdatedBy = parsedUserID
	updatedStudent.UpdatedOn = time.Now()

	return s.repo.Update(ctx, id, updatedStudent)
}

func (s *StudentService) DeleteStudent(ctx context.Context, id string) error {
	currentUserId := ctx.Value("userID").(string)
	existingStudent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("student not found")
	}
	if existingStudent.CreatedBy.String() != currentUserId {
		return errors.New("you are not allowed to delete this student")
	}
	return s.repo.Delete(ctx, id)
}
