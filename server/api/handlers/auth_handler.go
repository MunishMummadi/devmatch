package handlers

import (
	"log"
	"net/http"

	"github.com/MunishMummadi/devmatch/server/internal/services"
	"github.com/MunishMummadi/devmatch/server/internal/services/database"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	DBService    *database.DBService
	ClerkService *services.ClerkService
	// Add Clerk client or other services if needed directly
}

func NewAuthHandler(db *database.DBService, clerkService *services.ClerkService) *AuthHandler {
	return &AuthHandler{DBService: db, ClerkService: clerkService}
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
	// Use the base package's function to get claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}

	// Fetch the Clerk User object using the UserID from the session claims
	userProfile, err := h.ClerkService.GetUser(c.Request.Context(), claims.Subject)
	if err != nil {
		// Handle specific Clerk errors if necessary, otherwise generic error
		log.Printf("Error fetching user from Clerk API for UserID %s: %v", claims.Subject, err)
		// Check if it's a 'not found' error from Clerk if possible, otherwise assume internal error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile from Clerk"})
		return
	}

	// Check if the user profile exists in *our* database
	_, dbErr := h.DBService.GetUserProfileByClerkID(c.Request.Context(), userProfile.ID) // Use userProfile.ID from Clerk response
	if dbErr != nil {
		if dbErr == pgx.ErrNoRows {
			log.Printf("Profile not found in DB for Clerk User ID: %s\n", userProfile.ID)
			c.JSON(http.StatusNotFound, gin.H{
				"error":         "User profile not found in database. Please complete registration.",
				"clerkUserData": userProfile, // Optionally return Clerk data to aid frontend
			})
		} else {
			log.Printf("Error fetching profile from DB for Clerk User ID %s: %v\n", userProfile.ID, dbErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile from database"})
		}
		return
	}

	// If we reach here, Clerk user exists and profile exists in our DB
	// Return the profile from our DB (as it might have more/different info than Clerk's raw data)
	dbProfile, _ := h.DBService.GetUserProfileByClerkID(c.Request.Context(), userProfile.ID) // Fetch again to return
	c.JSON(http.StatusOK, dbProfile)                                                         // Return the profile from our DB
}
