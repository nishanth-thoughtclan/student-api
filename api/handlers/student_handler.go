package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/services"
)

// Add or Update Student Request type
type StudentRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

// Response objecgi
type Response struct {
	Message string `json:"message"`
}

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
		students, err := service.GetAllStudents(r.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(students) == 0 {
			json.NewEncoder(w).Encode([]models.Student{})
		} else {
			json.NewEncoder(w).Encode(students)
		}
	}
}

func GetStudentByIDHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		student, err := service.GetStudentByID(r.Context(), id)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(student)
	}
}

func CreateStudentHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var studentReq StudentRequest
		if err := json.NewDecoder(r.Body).Decode(&studentReq); err != nil {
			return
		}

		validate := validator.New()
		err := validate.Struct(studentReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		student := studentFromAddStudentRequest(studentReq)

		updatedStudent, saveErr := service.CreateStudent(r.Context(), student)
		if saveErr != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(updatedStudent); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func UpdateStudentHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		var studentReq StudentRequest
		if err := json.NewDecoder(r.Body).Decode(&studentReq); err != nil {
			return
		}

		validate := validator.New()
		err := validate.Struct(studentReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		student := studentFromAddStudentRequest(studentReq)
		updatedStudent, updateErr := service.UpdateStudent(r.Context(), id, student)
		if updateErr != nil {
			if updateErr.Error() == "student not found" {
				http.Error(w, updateErr.Error(), http.StatusNotFound)
			} else {
				http.Error(w, updateErr.Error(), http.StatusForbidden)
			}
			return
		}

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(updatedStudent); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func DeleteStudentHandler(service *services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		err := service.DeleteStudent(r.Context(), id)
		if err != nil {
			if err.Error() == "student not found" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusForbidden)
			}
			return
		}
		if err := json.NewEncoder(w).Encode(Response{Message: "Successfully Deleted"}); err != nil {
			panic(err)
		}
	}
}

func studentFromAddStudentRequest(r StudentRequest) models.Student {
	return models.Student{
		Name: r.Name,
		Age:  r.Age,
	}
}
