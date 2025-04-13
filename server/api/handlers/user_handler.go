package handlers

import (
	"log"
	"net/http"

	"server/api/middleware"             // <-- Adjust import path
	"server/internal/models"            // <-- Adjust import path
	"server/internal/services/database" // <-- Adjust import path

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DBService *database.DBService
	// Add GeminiService, GitHubService later if needed
}

func NewUserHandler(db *database.DBService) *UserHandler {
	return &UserHandler{DBService: db}
}

// GetUserProfileByID godoc
// @Summary Get user profile by their DB ID
// @Description Retrieves a user's public profile information by their internal database ID.
// @Tags Users
// @Produce json
// @Param id path string true "User Database ID (UUID or Int)"
// @Success 200 {object} models.User "Successfully retrieved user profile"
// @Failure 400 {object} gin.H "Invalid ID format"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserProfileByID(c *gin.Context) {
	userID := c.Param("id") // This is the DB ID (e.g., UUID), NOT Clerk ID

	// Basic validation (e.g., check if it's a valid UUID if using UUIDs)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID parameter is required"})
		return
	}

	// --- IMPORTANT ---
	// You need a DBService function like GetUserProfileByDBID(ctx, dbID)
	// The current GetUserProfileByClerkID won't work here directly.
	// For this example, we'll *assume* such a function exists or modify the Get function.
	// Let's *pretend* GetUserProfileByClerkID works with DB ID for demonstration simplicity,
	// BUT YOU SHOULD CREATE A SEPARATE FUNCTION GetUserProfileByDBID in db.go
	// For now, we will just return a placeholder error.

	// userProfile, err := h.DBService.GetUserProfileByDBID(c.Request.Context(), userID) // Correct approach
	// Placeholder:
	log.Printf("Attempting to get user profile by DB ID: %s (DB query function not implemented yet for this endpoint)", userID)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Fetching user by database ID not fully implemented yet"})

	/* // Example of how it *would* look with GetUserProfileByDBID:
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Error fetching profile for DB ID %s: %v\n", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		}
		return
	}
	c.JSON(http.StatusOK, userProfile)
	*/

	// TODO: Enrich with GitHub summary from Gemini API (requires GeminiService)
}

// CreateOrUpdateCurrentUserProfile godoc
// @Summary Create or update the current authenticated user's profile
// @Description Creates a new profile or updates existing profile data for the authenticated user. Uses Clerk ID from session.
// @Tags Users
// @Accept json
// @Produce json
// @Security ClerkAuth
// @Param profile body models.CreateUserProfileRequest true "Profile data to create or update"
// @Success 200 {object} models.User "Successfully updated profile"
// @Success 201 {object} models.User "Successfully created profile"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/profile [post] // Using a dedicated endpoint instead of /users/create
func (h *UserHandler) CreateOrUpdateCurrentUserProfile(c *gin.Context) {
	clerkUserID, exists := middleware.GetClerkUserID(c)
	if !exists {
		log.Println("Error: ClerkUserID not found in context in CreateOrUpdateCurrentUserProfile")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Cannot identify user"})
		return
	}

	var req models.CreateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for profile creation/update: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Map request data to the database model, injecting the Clerk User ID
	user := models.User{
		ClerkUserID: clerkUserID,
		Username:    req.Username,
		PictureURL:  req.PictureURL,
		Bio:         req.Bio,
		GitHubURL:   req.GitHubURL,
	}

	// Check if user already exists to return 200 (update) or 201 (create)
	_, err := h.DBService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	isUpdate := err == nil // If no error (found), it's an update

	createdOrUpdatedUser, err := h.DBService.CreateOrUpdateUserProfile(c.Request.Context(), user)
	if err != nil {
		log.Printf("Error saving profile for Clerk User ID %s: %v\n", clerkUserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user profile"})
		return
	}

	if isUpdate {
		c.JSON(http.StatusOK, createdOrUpdatedUser) // 200 OK for update
	} else {
		c.JSON(http.StatusCreated, createdOrUpdatedUser) // 201 Created for new profile
	}
}

// TODO: Implement PUT /users/:id (Edit existing profile info - requires careful auth checks)
// TODO: Implement GET /users/random (Get list of random users for swiping)
