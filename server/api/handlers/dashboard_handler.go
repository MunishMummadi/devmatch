package handlers

import (
	"net/http"

	"gin/api/middleware" // Needed for GetClerkUserID
	"gin/internal/services/database"

	"github.com/gin-gonic/gin"
	// Add other necessary imports like "gin/internal/models"
)

// DashboardHandler handles API requests related to the user dashboard (swiping, favorites).
type DashboardHandler struct {
	dbService *database.DBService
	// Add other services if needed (e.g., matching service)
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(db *database.DBService) *DashboardHandler {
	return &DashboardHandler{
		dbService: db,
	}
}

// GetSwipeCards fetches potential matches for the logged-in user.
// GET /dashboard/cards
func (h *DashboardHandler) GetSwipeCards(c *gin.Context) {
	userID, exists := middleware.GetClerkUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// TODO: Implement logic to fetch potential matches for userID from h.dbService
	// This might involve filtering out users already swiped, etc.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "TODO: Get swipe cards for user " + userID})
}

// LogSwipe records a swipe action.
// POST /dashboard/swipe
func (h *DashboardHandler) LogSwipe(c *gin.Context) {
	userID, exists := middleware.GetClerkUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// TODO: Bind request body (e.g., { "swiped_user_id": "...", "direction": "like|dislike" })
	// TODO: Implement logic to save swipe via h.dbService
	// TODO: Check for a match if direction is 'like'
	c.JSON(http.StatusNotImplemented, gin.H{"message": "TODO: Log swipe from user " + userID})
}

// ToggleFavorite adds or removes a user from the logged-in user's favorites.
// POST /dashboard/favorite
func (h *DashboardHandler) ToggleFavorite(c *gin.Context) {
	userID, exists := middleware.GetClerkUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// TODO: Bind request body (e.g., { "favorite_user_id": "...", "action": "add|remove" })
	// TODO: Implement logic to add/remove favorite via h.dbService
	c.JSON(http.StatusNotImplemented, gin.H{"message": "TODO: Toggle favorite for user " + userID})
}

// GetFavorites fetches the list of users favorited by the logged-in user.
// GET /dashboard/favorites
func (h *DashboardHandler) GetFavorites(c *gin.Context) {
	userID, exists := middleware.GetClerkUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// TODO: Implement logic to fetch favorite users for userID from h.dbService
	c.JSON(http.StatusNotImplemented, gin.H{"message": "TODO: Get favorites for user " + userID})
}