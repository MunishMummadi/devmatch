package routes

import (
	"log"
	"net/http"

	"github.com/MunishMummadi/devmatch/server/api/handlers"               // Use new module path
	"github.com/MunishMummadi/devmatch/server/internal/services"          // Use new module path
	"github.com/MunishMummadi/devmatch/server/internal/services/database" // Use new module path

	// Base clerk package needed for SessionClaimsFromContext in handlers
	// clerk "github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http" // HTTP package contains middleware
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// SetupRouter configures the Gin router with middleware and routes.
func SetupRouter(dbService *database.DBService, githubService *services.GitHubService, geminiService *services.GeminiService, clerkService *services.ClerkService) *gin.Engine {
	log.Println("Setting up Gin router...")
	// Set Gin mode (debug, release, test)
	ginMode := gin.DebugMode // Default or get from config
	if mode := gin.Mode(); mode != "" {
		ginMode = mode
	}
	gin.SetMode(ginMode)

	router := gin.Default() // Includes logger and recovery middleware

	// --- Middleware ---
	// CORS (Allow requests from your frontend) - Configure properly for production
	router.Use(cors.Default()) // Add default CORS configuration

	// Apply Clerk middleware - Use RequireHeaderAuthorization from http package
	router.Use(clerkGinMiddleware(clerkhttp.RequireHeaderAuthorization()))

	// Create services and handlers - Use passed-in services
	hub := handlers.NewHub(dbService)                           // Create the Hub instance
	// Pass both dbService and clerkService to NewAuthHandler
	authHandler := handlers.NewAuthHandler(dbService, clerkService)
	userHandler := handlers.NewUserHandler(dbService)           // Use passed-in dbService
	chatHandler := handlers.NewChatHandler(dbService, hub)      // Use passed-in dbService
	dashboardHandler := handlers.NewDashboardHandler(dbService) // Use passed-in dbService
	// Use passed-in githubService and geminiService
	githubHandler := handlers.NewGitHubHandler(githubService, geminiService)

	// --- Routes ---
	// Public Routes (e.g., health check, maybe docs)
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "UP"}) })
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Add Swagger UI route

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
		authGroup.GET("/user", authHandler.GetCurrentUserProfile)
		authGroup.POST("/logout", func(c *gin.Context) {
			// Backend can't easily invalidate Clerk session cookie (HttpOnly).
			// Needs coordination with frontend Clerk SDK (signOut()).
			// Placeholder - Consider if backend action is needed beyond frontend signout.
			c.JSON(http.StatusOK, gin.H{"message": "Logout initiated. Frontend should clear session."})
		})

		// Dashboard Routes (Swiping, Favorites) - Apply middleware to the group
		dashboardGroup := authGroup.Group("/dashboard")
		{
			dashboardGroup.GET("/cards", dashboardHandler.GetSwipeCards)      // Get potential matches
			dashboardGroup.POST("/swipe", dashboardHandler.LogSwipe)          // Log a swipe action
			dashboardGroup.POST("/favorite", dashboardHandler.ToggleFavorite) // Add/remove favorite
			dashboardGroup.GET("/favorites", dashboardHandler.GetFavorites)   // Get favorite users
		}

		// Chat Routes - Apply middleware to the group
		chatGroup := authGroup.Group("/chat")
		{
			// GET /auth/chat/conversations - User ID from auth middleware context
			chatGroup.GET("/conversations", chatHandler.GetConversations)
			// chatGroup.GET("/conversations/:userId", chatHandler.GetConversations) // Param redundant
			chatGroup.GET("/messages/:conversationId", chatHandler.GetMessages) // Get messages for a conversation
			chatGroup.POST("/message", chatHandler.SendMessage)                 // Send a message (REST)
			chatGroup.GET("/ws", chatHandler.HandleWebSocket)                   // WebSocket connection endpoint
		}

		// GitHub/Developer Tool Routes - Apply middleware to the group
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
		userGroup.POST("/profile", userHandler.CreateOrUpdateCurrentUserProfile)

		// Edit OWN profile - Requires Authentication and authorization check inside handler
		userGroup.PUT("/:id", userHandler.EditUserProfile) // Requires auth + check if user edits own profile

		// Get random users for swiping - Requires Authentication
		userGroup.GET("/random", userHandler.GetRandomUsers) // Requires auth to know who to exclude
	}

	log.Println("Router setup complete.")
	return router
}

// clerkGinMiddleware wraps the standard Clerk HTTP middleware for Gin.
func clerkGinMiddleware(stdMiddleware func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		dummyNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r // Update Gin context with request potentially modified by Clerk
			c.Next()    // Continue Gin chain
		})

		wrappedHandler := stdMiddleware(dummyNext)
		wrappedHandler.ServeHTTP(c.Writer, c.Request)

		// Abort Gin chain if Clerk middleware wrote response (e.g., 401/403)
		if c.Writer.Written() {
			c.Abort()
		}
	}
}
