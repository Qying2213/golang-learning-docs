package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qinyang/taskmanager/config"
	"github.com/qinyang/taskmanager/database"
	"github.com/qinyang/taskmanager/handlers"
	"github.com/qinyang/taskmanager/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg.Database.DSN); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Gin
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	taskHandler := handlers.NewTaskHandler()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// User profile
			protected.GET("/profile", authHandler.GetProfile)

			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", taskHandler.GetTasks)
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
			}
		}
	}

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddr)

	server := &http.Server{
		Addr:           serverAddr,
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
