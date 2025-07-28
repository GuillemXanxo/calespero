package api

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"calespero/internal/core/services"
	"calespero/internal/handlers"
	"calespero/internal/repositories/postgres"
	"calespero/pkg/auth"

	_ "github.com/lib/pq"
)

func Run() {
	// Load templates
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// Initialize DB connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories and services
	userRepo := postgres.NewUserRepository(db)
	jwtManager := auth.NewJWTManager(
		os.Getenv("JWT_SECRET_KEY"),
		24*time.Hour,
	)
	userSvc := services.NewUserService(userRepo, jwtManager)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(templates, userSvc)

	// Setup routes
	http.HandleFunc("/", userHandler.HandleLogin)
	http.HandleFunc("/login", userHandler.HandleLogin)
	http.HandleFunc("/new_user", userHandler.HandleNewUser)
	http.HandleFunc("/start", userHandler.HandleStart)

	// Start server
	log.Printf("Server starting on port 3500...")
	log.Fatal(http.ListenAndServe(":3500", nil))
}
