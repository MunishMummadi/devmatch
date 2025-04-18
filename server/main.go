package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin/api/routes"                 // Corrected import path
	"gin/internal/config"            // Corrected import path
	"gin/internal/services"          // Added services import
	"gin/internal/services/database" // Corrected import path

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// @title Your Project API
// @version 1.0
// @description This is the API server for the Your Project application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ClerkAuth
// @in header
// @name Authorization
// @description Clerk JWT token (include 'Bearer ' prefix)

func main() {
	log.Println("Starting server...")

	// Load Configuration
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: Port=%s, GinMode=%s\n", cfg.Port, cfg.GinMode)

	// Initialize Clerk Client
	clerkClient, err := clerk.NewClient(cfg.ClerkSecretKey)
	if err != nil {
		log.Fatalf("Failed to create Clerk client: %v", err)
	}
	log.Println("Clerk client initialized successfully.")

	// Initialize Database Connection Pool
	dbPool, err := database.ConnectDB(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer func() {
		log.Println("Closing database connection pool...")
		dbPool.Close()
	}()
	log.Println("Database connection pool initialized successfully.")

	// Initialize Services
	githubService := services.NewGitHubService(cfg)
	geminiService, err := services.NewGeminiService(cfg)
	if err != nil {
		// Decide how to handle Gemini init failure - fatal or just log?
		// For now, log and continue, but the service will be nil/unusable.
		// log.Fatalf("Failed to initialize Gemini service: %v", err) // Option: Fatal
		log.Printf("Warning: Failed to initialize Gemini service: %v. Gemini features will be disabled.", err)
		// Ensure geminiService is nil if error occurred (already is due to return nil)
	} else {
		// Ensure Gemini client is closed on shutdown
		defer func() {
			log.Println("Closing Gemini client...")
			geminiService.Close()
		}()
	}
	clerkService := services.NewClerkService(clerkClient, cfg)
	log.Println("Application services initialized.")

	// Setup Gin Router
	router := routes.SetupRouter(dbPool, clerkClient, githubService, geminiService, clerkService)
	log.Println("Gin router setup complete.")

	// Setup HTTP Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
		// Optional: Add timeouts for production readiness
		// ReadTimeout:  15 * time.Second,
		// WriteTimeout: 15 * time.Second,
		// IdleTimeout:  60 * time.Second,
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		log.Printf("Server listening on port %s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful Shutdown Handling
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the requests it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting.")
}
