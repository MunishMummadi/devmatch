package handlers

import (
	"log"
	"net/http"

	"server/api/middleware"             // <-- Adjust import path
	"server/internal/services/database" // <-- Adjust import path

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	DBService *database.DBService
	// Add Clerk client or other services if needed directly
}

func NewAuthHandler(db *database.DBService) *AuthHandler {
	return &AuthHandler{DBService: db}
}

// GetCurrentUserProfile godoc
// @Summary Get current authenticated user's profile
// @Description Retrieves the profile data for the user authenticated via Clerk.
// @Tags Auth
// @Produce json
// @Security ClerkAuth
// @Success 200 {object} models.User "Successfully retrieved user profile"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "User profile not found in DB"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /auth/user [get]
func (h *AuthHandler) GetCurrentUserProfile(c *gin.Context) {
	clerkUserID, exists := middleware.GetClerkUserID(c)
	if !exists {
		// This shouldn't happen if middleware is applied correctly, but good practice to check
		log.Println("Error: ClerkUserID not found in context in GetCurrentUserProfile")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not identify authenticated user"})
		return
	}

	userProfile, err := h.DBService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Profile not found for Clerk User ID: %s\n", clerkUserID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found. Please create one."})
		} else {
			log.Printf("Error fetching profile for Clerk User ID %s: %v\n", clerkUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		}
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

// TODO: Implement /auth/github/login (Redirect to Clerk's GitHub handler)
// TODO: Implement /auth/github/callback (Handled by Clerk frontend components usually, backend might just need to ensure session is created)
// TODO: Implement /auth/logout (Needs coordination with Clerk frontend SDK for clearing cookies/session)
