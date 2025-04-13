package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/MunishMummadi/devmatch/server/internal/models"
	"github.com/MunishMummadi/devmatch/server/internal/services/database"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles API requests related to the user dashboard (swiping, favorites).
type DashboardHandler struct {
	dbService *database.DBService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(db *database.DBService) *DashboardHandler {
	return &DashboardHandler{
		dbService: db,
	}
}

// GetSwipeCards godoc
// @Summary Get swipe cards
// @Description Fetches a list of potential user profiles for the logged-in user to swipe on.
// @Tags Dashboard
// @Produce json
// @Security ClerkAuth
// @Param limit query int false "Number of cards to return" default(10)
// @Success 200 {array} models.User "List of user profiles"
// @Failure 400 {object} gin.H "Invalid limit parameter"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Authenticated user profile not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /dashboard/cards [get]
func (h *DashboardHandler) GetSwipeCards(c *gin.Context) {
	// Get User ID from Clerk claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

	// Fetch current user's DB ID
	currentUser, err := h.dbService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error: Authenticated user %s not found in DB for GetSwipeCards", clerkUserID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user profile not found"})
		} else {
			log.Printf("Error fetching current user profile %s in GetSwipeCards: %v\n", clerkUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user profile"})
		}
		return
	}

	// Get Limit Parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' query parameter"})
		return
	}
	if limit > 50 { // Optional server-side cap
		limit = 50
	}

	// Fetch random users (excluding self and already swiped)
	// Filtering logic (excluding self and swiped) is now handled within dbService.GetRandomUsers
	swipeCards, err := h.dbService.GetRandomUsers(c.Request.Context(), currentUser.ID, limit)
	if err != nil {
		log.Printf("Error getting random users for swipe cards (excluding %s): %v\n", currentUser.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve swipe cards"})
		return
	}

	c.JSON(http.StatusOK, swipeCards)
}

type SwipeRequest struct {
	SwipedUserID string `json:"swiped_user_id" binding:"required"`
	Direction    string `json:"direction" binding:"required,oneof=like dislike"` // Use string directly
}

// LogSwipe godoc
// @Summary Log a swipe action
// @Description Records a user's swipe action (like/dislike) on another user and checks for matches.
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security ClerkAuth
// @Param swipe body SwipeRequest true "Swipe action details"
// @Success 200 {object} gin.H "Swipe logged successfully"
// @Success 201 {object} gin.H{match=bool,conversation_id=string} "Match found!" // Include conversation ID on match
// @Failure 400 {object} gin.H "Invalid request body or swipe direction"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Current user or swiped user not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /dashboard/swipe [post]
func (h *DashboardHandler) LogSwipe(c *gin.Context) {
	// Get User ID from Clerk claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

	// Fetch current user's DB ID
	currentUser, err := h.dbService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		// Handle error similar to GetSwipeCards
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user profile not found"})
		} else {
			log.Printf("Error fetching current user profile %s in LogSwipe: %v\n", clerkUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user profile"})
		}
		return
	}
	currentUserID := currentUser.ID // The DB ID

	// Bind request body
	var req SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding swipe request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Ensure user isn't swiping on themselves (shouldn't happen with GetSwipeCards filtering, but belt-and-suspenders)
	if req.SwipedUserID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot swipe on yourself"})
		return
	}

	// Save the swipe
	swipe := models.Swipe{
		SwiperID:  currentUserID,                        // Correct field name
		SwipedID:  req.SwipedUserID,                     // Correct field name
		Direction: models.SwipeDirection(req.Direction), // Convert string to models.SwipeDirection
	}
	err = h.dbService.SaveSwipe(c.Request.Context(), swipe)
	if err != nil {
		// Could be foreign key constraint if swiped_user_id doesn't exist
		log.Printf("Error saving swipe from %s to %s (%s): %v\n", currentUserID, req.SwipedUserID, req.Direction, err)
		// Check for specific DB errors if needed, otherwise generic 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save swipe"})
		return
	}

	// Check for match only if it was a 'like'
	if swipe.Direction == models.SwipeLike { // Compare with the constant
		isMatch, err := h.dbService.CheckForMatch(c.Request.Context(), currentUserID, req.SwipedUserID)
		if err != nil {
			log.Printf("Error checking for match between %s and %s: %v\n", currentUserID, req.SwipedUserID, err)
			// Continue without match confirmation, but log the error
			c.JSON(http.StatusOK, gin.H{"message": "Swipe logged successfully (match check failed)"}) // Or maybe 500?
			return
		}

		if isMatch {
			log.Printf("Match found between User %s and User %s!\n", currentUserID, req.SwipedUserID)
			// Create conversation if it doesn't exist
			participantIDs := []string{currentUserID, req.SwipedUserID}
			conversation, err := h.dbService.CreateConversation(c.Request.Context(), participantIDs)
			if err != nil {
				log.Printf("Error creating/finding conversation between %s and %s: %v", currentUserID, req.SwipedUserID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate conversation"})

			}
			log.Printf("Created conversation ID %s for match.\n", conversation.ID) // Assuming ID is string or int
			c.JSON(http.StatusCreated, gin.H{
				"message":         "Match found!",
				"match":           true,
				"conversation_id": conversation.ID, // Remove strconv.Itoa
			})
			return
		}
	}

	// If not a like, or like but no match
	c.JSON(http.StatusOK, gin.H{"message": "Swipe logged successfully", "match": false})
}

type FavoriteRequest struct {
	FavoriteUserID string `json:"favorite_user_id" binding:"required"`
}

// ToggleFavorite adds or removes a user from the logged-in user's favorites.
// POST /dashboard/favorite
// @Summary Toggle favorite status
// @Description Adds or removes a user from the logged-in user's favorites list.
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security ClerkAuth
// @Param favorite body FavoriteRequest true "User ID to toggle favorite status for"
// @Success 200 {object} gin.H "Favorite status updated (added/removed)"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Current user or target user not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /dashboard/favorite [post]
func (h *DashboardHandler) ToggleFavorite(c *gin.Context) {
	// Get User ID from Clerk claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

	// Fetch current user's DB ID
	currentUser, err := h.dbService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user profile not found"})
		} else {
			log.Printf("Error fetching current user profile %s in ToggleFavorite: %v\n", clerkUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user profile"})
		}
		return
	}
	currentUserID := currentUser.ID

	// Bind request body
	var req FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding favorite request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Prevent favoriting self
	if req.FavoriteUserID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot favorite yourself"})
		return
	}

	// Check current favorite status
	isCurrentlyFavorite, err := h.dbService.IsFavorite(c.Request.Context(), currentUserID, req.FavoriteUserID)
	if err != nil {
		log.Printf("Error checking favorite status for user %s by user %s: %v\n", req.FavoriteUserID, currentUserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite status"})
		return
	}

	var actionMessage string
	if isCurrentlyFavorite {
		// Remove favorite
		err = h.dbService.RemoveFavorite(c.Request.Context(), currentUserID, req.FavoriteUserID)
		if err != nil {
			log.Printf("Error removing favorite %s for user %s: %v\n", req.FavoriteUserID, currentUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
			return
		}
		actionMessage = "User removed from favorites"
		log.Printf("User %s removed favorite %s\n", currentUserID, req.FavoriteUserID)
	} else {
		// Add favorite
		err = h.dbService.AddFavorite(c.Request.Context(), currentUserID, req.FavoriteUserID)
		if err != nil {
			// Could be foreign key constraint if favorite_user_id doesn't exist
			log.Printf("Error adding favorite %s for user %s: %v\n", req.FavoriteUserID, currentUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"}) // Potentially 404 if user doesn't exist? Depends on DB constraints
			return
		}
		actionMessage = "User added to favorites"
		log.Printf("User %s added favorite %s\n", currentUserID, req.FavoriteUserID)
	}

	c.JSON(http.StatusOK, gin.H{"message": actionMessage})
}

// GetFavorites fetches the list of users favorited by the logged-in user.
// GET /dashboard/favorites
// @Summary Get favorite users
// @Description Retrieves the list of users that the logged-in user has favorited.
// @Tags Dashboard
// @Produce json
// @Security ClerkAuth
// @Success 200 {array} models.User "List of favorite users"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Current user profile not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /dashboard/favorites [get]
func (h *DashboardHandler) GetFavorites(c *gin.Context) {
	// Get User ID from Clerk claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

	// Fetch current user's DB ID
	currentUser, err := h.dbService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user profile not found"})
		} else {
			log.Printf("Error fetching current user profile %s in GetFavorites: %v\n", clerkUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user profile"})
		}
		return
	}
	currentUserID := currentUser.ID

	// Fetch favorites from DB
	favorites, err := h.dbService.GetFavoritesByUserID(c.Request.Context(), currentUserID)
	if err != nil {
		// sql.ErrNoRows is not expected here, as an empty favorite list is valid.
		// Log the specific error for debugging.
		log.Printf("Error fetching favorites for user %s: %v\n", currentUserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve favorites"})
		return
	}

	// Return the list (even if empty)
	c.JSON(http.StatusOK, favorites)
}
