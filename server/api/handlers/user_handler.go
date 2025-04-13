package handlers

import (
	"database/sql" // Added for sql.ErrNoRows
	"errors"       // Added for errors.Is
	"log"
	"net/http"
	"strconv" // Added for parsing limit

	"github.com/MunishMummadi/devmatch/server/internal/models"            // Use new module path
	"github.com/MunishMummadi/devmatch/server/internal/services/database" // Use new module path
	clerk "github.com/clerk/clerk-sdk-go/v2"                              // Use new module path

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
// @Param id path string true "User Database ID (UUID or other string format)"
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

	userProfile, err := h.DBService.GetUserProfileByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Error fetching profile for DB ID %s: %v\n", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		}
		return
	}

	// TODO: Enrich with GitHub summary from Gemini API (requires GeminiService) - This can be done here after fetching the basic profile

	c.JSON(http.StatusOK, userProfile)
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
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		log.Println("Error: Clerk session claims not found in context in CreateOrUpdateCurrentUserProfile")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

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

// EditUserProfile godoc
// @Summary Edit the current authenticated user's profile
// @Description Updates specific fields of the authenticated user's profile information. Only allows editing own profile.
// @Tags Users
// @Accept json
// @Produce json
// @Security ClerkAuth
// @Param id path string true "User Database ID (Must match authenticated user)"
// @Param profile body models.CreateUserProfileRequest true "Profile data fields to update"
// @Success 200 {object} models.User "Successfully updated user profile"
// @Failure 400 {object} gin.H "Invalid request body or ID format"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden - Cannot edit another user's profile"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) EditUserProfile(c *gin.Context) {
	// 1. Get Authenticated User's Clerk ID
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		log.Println("Error: Clerk session claims not found in context in EditUserProfile")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

	// 2. Get Target User DB ID from Path
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID parameter is required"})
		return
	}

	// 3. Fetch Target User Profile
	targetUserProfile, err := h.DBService.GetUserProfileByID(c.Request.Context(), targetUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Error fetching target profile for DB ID %s in EditUserProfile: %v\n", targetUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		}
		return
	}

	// 4. Authorization Check: Clerk ID from token must match Clerk ID of the profile being edited
	if targetUserProfile.ClerkUserID != clerkUserID {
		log.Printf("Forbidden attempt: User %s tried to edit profile of user %s (ClerkID: %s)\n", clerkUserID, targetUserID, targetUserProfile.ClerkUserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only edit your own profile"})
		return
	}

	// 5. Bind Request Body
	var req models.CreateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for profile update: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 6. Call DBService to Edit
	// Note: We pass targetUserID (the DB ID) here, as EditUserProfile works on DB ID
	updatedUser, err := h.DBService.EditUserProfile(c.Request.Context(), targetUserID, req)
	if err != nil {
		// EditUserProfile might return ErrNoRows if the user disappeared between check and update, unlikely but handle
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found during update"})
		} else {
			log.Printf("Error updating profile for User DB ID %s: %v\n", targetUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		}
		return
	}

	// 7. Return Success
	c.JSON(http.StatusOK, updatedUser)
}

// GetRandomUsers godoc
// @Summary Get random user profiles for swiping
// @Description Retrieves a list of random user profiles, excluding the current user.
// @Tags Users
// @Produce json
// @Security ClerkAuth
// @Param limit query int false "Number of users to return" default(10)
// @Success 200 {array} models.User "Successfully retrieved random users"
// @Failure 400 {object} gin.H "Invalid limit parameter"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Authenticated user profile not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/random [get]
func (h *UserHandler) GetRandomUsers(c *gin.Context) {
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		log.Println("Error: Clerk session claims not found in context in GetRandomUsers")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	clerkUserID := claims.Subject

	// 2. Fetch Authenticated User's DB ID (needed to exclude them)
	currentUserProfile, err := h.DBService.GetUserProfileByClerkID(c.Request.Context(), clerkUserID)
	if err != nil {
		// This shouldn't happen if the user is authenticated, but handle defensively
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error: Authenticated user %s not found in DB for GetRandomUsers", clerkUserID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user profile not found"})
		} else {
			log.Printf("Error fetching current user profile %s in GetRandomUsers: %v\n", clerkUserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current user profile"})
		}
		return
	}

	// 3. Get Limit Parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' query parameter"})
		return
	}
	// Optional: Add a server-side maximum limit
	if limit > 50 {
		limit = 50
	}

	// 4. Call DBService
	randomUsers, err := h.DBService.GetRandomUsers(c.Request.Context(), currentUserProfile.ID, limit)
	if err != nil {
		log.Printf("Error getting random users (excluding %s): %v\n", currentUserProfile.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve random users"})
		return
	}

	// 5. Return Success (even if the list is empty)
	c.JSON(http.StatusOK, randomUsers)
}
