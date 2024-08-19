package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/nishanth-thoughtclan/student-api/api/handlers"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/config"
	_ "github.com/nishanth-thoughtclan/student-api/docs"
	"github.com/nishanth-thoughtclan/student-api/middlewares"
	middleware "github.com/nishanth-thoughtclan/student-api/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// loading .env vars
	cfg := config.LoadConfig()
	// connecting to the db
	db, err := sql.Open("mysql", cfg.DBUser+":"+cfg.DBPassword+"@tcp("+cfg.DBHost+")/"+cfg.DBName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	studentRepo := repositories.NewStudentRepository(db)
	userRepo := repositories.NewUserRepository(db)
	studentService := services.NewStudentService(studentRepo)
	authService := services.NewAuthService(userRepo)
	// register routes
	router := mux.NewRouter()
	registerRoutesAndMiddlewares(router, authService, studentService)
	// starting the server, default 8080
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
	log.Printf("Server started on port %s", port)
	defer db.Close()
}

func registerRoutesAndMiddlewares(router *mux.Router, authService *services.AuthService, studentService *services.StudentService) {
	router.Use(middlewares.LoggingMiddleware)
	router.Use(middleware.JSONMiddleware)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/users/login", handlers.AuthHandler(authService)).Methods("POST")
	router.HandleFunc("/users/signup", handlers.SingUpHandler(authService)).Methods("POST")
	router.Handle("/students", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.GetStudentsHandler(studentService)))).Methods("GET")
	router.Handle("/students/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.GetStudentByIDHandler(studentService)))).Methods("GET")
	router.Handle("/students", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.CreateStudentHandler(studentService)))).Methods("POST")
	router.Handle("/students/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.UpdateStudentHandler(studentService)))).Methods("PUT")
	router.Handle("/students/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.DeleteStudentHandler(studentService)))).Methods("DELETE")
}
