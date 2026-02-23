package main

import (
	"log"
	"time"

	"recruitment-platform/config"
	"recruitment-platform/handlers"
	"recruitment-platform/middleware"
	"recruitment-platform/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	log.Println("Connected to database")

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// AI Service
	aiService := &services.AIService{
		DB:          db,
		WhisperURL:  cfg.WhisperURL,
		OllamaURL:   cfg.OllamaURL,
		OllamaModel: cfg.OllamaModel,
	}
	log.Printf("AI Service configured: Whisper=%s, Ollama=%s (model: %s)", cfg.WhisperURL, cfg.OllamaURL, cfg.OllamaModel)

	// Handlers
	candidateHandler := &handlers.CandidateHandler{DB: db, AIService: aiService}
	questionnaireHandler := &handlers.QuestionnaireHandler{DB: db}
	adminHandler := &handlers.AdminHandler{DB: db, JWTSecret: cfg.JWTSecret, AIService: aiService}

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/candidates", candidateHandler.Register)
		api.GET("/questionnaire", questionnaireHandler.GetActive)
		api.POST("/candidates/:id/answers", candidateHandler.SubmitAnswers)
		api.POST("/candidates/:id/audio", candidateHandler.UploadAudio)
	}

	// Admin routes (JWT protected)
	admin := api.Group("/admin")
	admin.POST("/login", adminHandler.Login)

	adminProtected := admin.Group("")
	adminProtected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		adminProtected.GET("/candidates", adminHandler.ListCandidates)
		adminProtected.GET("/candidates/:id", adminHandler.GetCandidate)
		adminProtected.GET("/candidates/:id/audio", adminHandler.GetCandidateAudio)
		adminProtected.POST("/candidates/:id/reanalyze", adminHandler.ReanalyzeCandidate)
		adminProtected.GET("/questionnaires", adminHandler.ListQuestionnaires)
		adminProtected.GET("/questionnaires/:id", adminHandler.GetQuestionnaire)
		adminProtected.POST("/questionnaires", adminHandler.CreateQuestionnaire)
	}

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
