package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

const (
	// Context key to store Clerk claims/session information
	ClerkSessionContextKey = "clerk_session"
	ClerkUserIDContextKey  = "clerk_user_id"
)

// ClerkMiddleware creates Gin middleware for Clerk authentication.
func ClerkMiddleware(clerkClient clerk.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := clerkClient.AuthenticateRequest(c.Request)
		if err != nil {
			// Log different types of errors
			log.Printf("Clerk authentication error: %v\n", err)

			// Handle specific Clerk error types if needed, otherwise generic unauthorized
			// Examples: errors.Is(err, clerk.ErrNoSessionToken), errors.Is(err, clerk.ErrInvalidSessionToken) etc.
			if strings.Contains(err.Error(), "no session token") || strings.Contains(err.Error(), "invalid session token") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing session token"})
			} else if strings.Contains(err.Error(), "expired") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Session token expired"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"}) // Or 401 depending on policy
			}
			return
		}

		if session == nil || session.Claims == nil {
			log.Println("Clerk authentication failed: No active session or claims found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No active session"})
			return
		}

		// Session is valid, enrich context
		c.Set(ClerkSessionContextKey, session)
		c.Set(ClerkUserIDContextKey, session.Claims.Subject) // Subject usually holds the User ID

		// Continue down the chain
		c.Next()
	}
}

// GetClerkUserID retrieves the Clerk User ID from the Gin context.
// Should be called *after* ClerkMiddleware has run successfully.
func GetClerkUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(ClerkUserIDContextKey)
	if !exists {
		return "", false
	}
	idStr, ok := userID.(string)
	return idStr, ok
}