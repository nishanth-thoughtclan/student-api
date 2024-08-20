package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nishanth-thoughtclan/student-api/api/handlers"
	"github.com/nishanth-thoughtclan/student-api/api/repositories"
	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/config"
	middleware "github.com/nishanth-thoughtclan/student-api/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/nishanth-thoughtclan/student-api/docs"
)

// @title           Student API
// @version         1.0
// @description     This is a sample server for a Student Management API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Open the log file explicitly
	logFile, err := os.OpenFile("D:\\student-api\\app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer logFile.Close()

	// Set the output of the standard logger to the log file
	log.SetOutput(logFile)
	// loading .env vars
	cfg := config.LoadConfig()
	// connecting to the db
	db, err := sql.Open("mysql", cfg.DBUser+":"+cfg.DBPassword+"@tcp("+cfg.DBHost+")/"+cfg.DBName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	studentRepo := repositories.NewStudentRepository(db)
	userRepo := repositories.NewUserRepository(db)
	studentService := services.NewStudentService(studentRepo)
	authService := services.NewAuthService(userRepo)
	// register routes
	router := mux.NewRouter()
	registerRoutesAndMiddlewares(router, db, authService, studentService)
	// starting the server, default 8080
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
		IdleTimeout:  time.Second * 180,
		Handler:      router,
	}

	// gracefully handle server shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s", err)
		}
	}()
	log.Printf("Server started on port %s", port)

	// wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}

func registerRoutesAndMiddlewares(router *mux.Router, db *sql.DB, authService *services.AuthService, studentService *services.StudentService) {
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.JSONMiddleware)
	router.HandleFunc("/health", handlers.ServiceHealthCheck).Methods("GET")
	router.HandleFunc("/ready", handlers.PingHandler(db)).Methods("GET")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/api/v1/users/login", handlers.AuthHandler(authService)).Methods("POST")
	router.HandleFunc("/api/v1/users/signup", handlers.SingUpHandler(authService)).Methods("POST")
	router.Handle("/api/v1/students", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.GetStudentsHandler(studentService)))).Methods("GET")
	router.Handle("/api/v1/students/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.GetStudentByIDHandler(studentService)))).Methods("GET")
	router.Handle("/api/v1/students", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.CreateStudentHandler(studentService)))).Methods("POST")
	router.Handle("/api/v1/students/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.UpdateStudentHandler(studentService)))).Methods("PUT")
	router.Handle("/api/v1/students/{id}", middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.DeleteStudentHandler(studentService)))).Methods("DELETE")
}
