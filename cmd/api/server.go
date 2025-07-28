package api

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"calespero/internal/config"
	"calespero/internal/core/services"
	"calespero/internal/handlers"
	"calespero/internal/middleware"
	"calespero/internal/repositories/postgres"
	"calespero/pkg/auth"
	"calespero/pkg/logger"
)

func Run() {
	// Initialize logger
	if err := logger.Initialize(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Starting application...")

	// Load templates
	templates := template.Must(template.ParseGlob("templates/*.html"))
	logger.Info("Templates loaded successfully")

	// Initialize DB connection
	dbConfig := config.NewDBConfigFromEnv()
	db, err := dbConfig.Connect()
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Database connection established")

	// Initialize repositories and services
	userRepo := postgres.NewUserRepository(db)
	jwtManager := auth.NewJWTManager(
		os.Getenv("JWT_SECRET_KEY"),
		24*time.Hour,
	)
	userSvc := services.NewUserService(userRepo, jwtManager)
	logger.Info("Services initialized")

	// Initialize handlers
	userHandler := handlers.NewUserHandler(templates, userSvc)

	// Setup routes with logging middleware
	http.HandleFunc("/", middleware.LoggingMiddleware(userHandler.HandleLogin))
	http.HandleFunc("/login", middleware.LoggingMiddleware(userHandler.HandleLogin))
	http.HandleFunc("/new_user", middleware.LoggingMiddleware(userHandler.HandleNewUser))
	http.HandleFunc("/start", middleware.LoggingMiddleware(userHandler.HandleStart))

	logger.Info("Routes configured, starting server on port 3500...")

	// Start server
	if err := http.ListenAndServe(":3500", nil); err != nil {
		logger.Error("Server failed to start: %v", err)
		os.Exit(1)
	}
}
