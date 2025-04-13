package routes

import (
	"database/sql"
	"log"
	"net/http"

	"gin/api/handlers"               // Corrected import path
	"gin/api/middleware"             // Corrected import path
	"gin/internal/services"          // Added services import
	"gin/internal/services/database" // Corrected import path

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

func SetupRouter(dbPool *sql.DB, clerkClient clerk.Client, githubService *services.GitHubService, geminiService *services.GeminiService, clerkService *services.ClerkService) *gin.Engine {
	// Set Gin mode (debug, release, test)
	ginMode := gin.DebugMode // Default or get from config
	if mode := gin.Mode(); mode != "" {
		ginMode = mode
	}
	gin.SetMode(ginMode)

	router := gin.Default() // Includes logger and recovery middleware

	// --- Middleware ---
	// CORS (Allow requests from your frontend) - Configure properly for production
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin (adjust for prod)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Create services and handlers
	dbService := database.NewDBService(dbPool)
	authHandler := handlers.NewAuthHandler(dbService)
	userHandler := handlers.NewUserHandler(dbService)
	chatHandler := handlers.NewChatHandler(dbService)
	dashboardHandler := handlers.NewDashboardHandler(dbService)
	githubHandler := handlers.NewGitHubHandler(githubService, geminiService)

	// Clerk Authentication Middleware Instance
	authMiddleware := middleware.ClerkMiddleware(clerkClient)

	// --- Routes ---
	// Public Routes (e.g., health check, maybe docs)
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "UP"}) })

	// Optional: Swagger Docs endpoint (if using Swaggo)
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authentication Routes
	authGroup := router.Group("/auth")
	{
		// These often involve redirects handled by Clerk's frontend SDKs / Hosted Pages
		// Placeholder endpoints - Actual implementation depends heavily on Clerk setup
		authGroup.GET("/github/login", func(c *gin.Context) {
			// Typically redirect to Clerk's hosted auth page or use frontend SDK
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Redirect to Clerk/GitHub handled by frontend or Clerk hosted pages"})
		})
		authGroup.GET("/github/callback", func(c *gin.Context) {
			// Clerk handles the callback and sets cookies/tokens
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Callback processed by Clerk; frontend should redirect"})
		})

		// Requires Authentication via Clerk session
		authGroup.GET("/user", authMiddleware, authHandler.GetCurrentUserProfile)
		authGroup.POST("/logout", authMiddleware, func(c *gin.Context) {
			// Backend can't easily invalidate Clerk session cookie (HttpOnly).
			// Needs coordination with frontend Clerk SDK (signOut()).
			c.JSON(http.StatusOK, gin.H{"message": "Logout initiated. Frontend should clear session."})
		})

		// Dashboard Routes (Swiping, Favorites)
		dashboardGroup := authGroup.Group("/dashboard")
		{
			dashboardGroup.GET("/cards", dashboardHandler.GetSwipeCards)      // Get potential matches
			dashboardGroup.POST("/swipe", dashboardHandler.LogSwipe)          // Log a swipe action
			dashboardGroup.POST("/favorite", dashboardHandler.ToggleFavorite) // Add/remove favorite
			dashboardGroup.GET("/favorites", dashboardHandler.GetFavorites)   // Get favorite users
		}

		// Chat Routes
		chatGroup := authGroup.Group("/chat")
		{
			chatGroup.GET("/conversations/:userId", chatHandler.GetConversations) // Get user's conversations (param might be redundant)
			chatGroup.GET("/messages/:conversationId", chatHandler.GetMessages)   // Get messages for a conversation
			chatGroup.POST("/message", chatHandler.SendMessage)                   // Send a message
		}

		// GitHub/Developer Tool Routes
		githubGroup := authGroup.Group("/github")
		{
			githubGroup.GET("/:username/data", githubHandler.GetGitHubData) // Fetch raw GitHub data
			githubGroup.POST("/summary", githubHandler.SummarizeGitHubData) // Summarize provided data (e.g., from GitHub)
		}
	}

	// User Profile Routes
	userGroup := router.Group("/users")
	{
		// Get public profile - potentially doesn't require auth depending on your rules
		userGroup.GET("/:id", userHandler.GetUserProfileByID) // :id is DB ID

		// Create or Update OWN profile - Requires Authentication
		userGroup.POST("/profile", authMiddleware, userHandler.CreateOrUpdateCurrentUserProfile)

		// TODO: userGroup.PUT("/:id", authMiddleware, userHandler.EditUserProfile) // Requires auth + check if user edits own profile
		// TODO: userGroup.GET("/random", authMiddleware, userHandler.GetRandomUsers) // Requires auth?
	}

	log.Println("Router setup complete.")
	return router
}
