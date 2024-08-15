package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/utils"

	"github.com/nishanth-thoughtclan/student-api/api/models"

	"github.com/gorilla/mux"
)

// GetStudentsHandler retrieves all students
// @Summary Get all students
// @Description Retrieves a list of all students
// @Tags Students
// @Produce json
// @Success 200 {array} Student
// @Failure 401 {object} ErrorResponse
// @Router /students [get]
func GetStudentsHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := service.GetAllStudents()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(students)
	}
}

func GetStudentByIDHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		student, err := service.GetStudentByID(id)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(student)
	}
}

func CreateStudentHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		json.NewDecoder(r.Body).Decode(&student)
		claims, _ := utils.ValidateJWTToken(r.Header.Get("Authorization"))
		student.CreatedBy = claims["sub"].(string)
		student.CreatedOn = time.Now()
		err := service.CreateStudent(student)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateStudentHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		var student models.Student
		json.NewDecoder(r.Body).Decode(&student)
		claims, _ := utils.ValidateJWTToken(r.Header.Get("Authorization"))
		student.CreatedBy = claims["sub"].(string)
		student.UpdatedOn = time.Now()
		err := service.UpdateStudent(id, student)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteStudentHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		err := service.DeleteStudent(id)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
