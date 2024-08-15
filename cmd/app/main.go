package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/nishanth-thoughtclan/student-api/api/handlers"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/config"
	"github.com/nishanth-thoughtclan/student-api/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("mysql", cfg.DBUser+":"+cfg.DBPassword+"@tcp("+cfg.DBHost+")/"+cfg.DBName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	studentRepo := repositories.NewStudentRepository(db)
	userRepo := repositories.NewUserRepository(db)
	studentService := services.NewStudentService(studentRepo)
	authService := services.NewAuthService(userRepo)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.HandleFunc("/auth", handlers.AuthHandler(cfg, authService)).Methods("POST")
	r.Handle("/students", middleware.FirebaseAuthMiddleware(http.HandlerFunc(handlers.GetStudentsHandler(studentService)))).Methods("GET")
	r.Handle("/students/{id}", middleware.FirebaseAuthMiddleware(http.HandlerFunc(handlers.GetStudentByIDHandler(studentService)))).Methods("GET")
	r.Handle("/students", middleware.FirebaseAuthMiddleware(http.HandlerFunc(handlers.CreateStudentHandler(studentService)))).Methods("POST")
	r.Handle("/students/{id}", middleware.FirebaseAuthMiddleware(http.HandlerFunc(handlers.UpdateStudentHandler(studentService)))).Methods("PUT")
	r.Handle("/students/{id}", middleware.FirebaseAuthMiddleware(http.HandlerFunc(handlers.DeleteStudentHandler(studentService)))).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
