package middleware

import (
	"log"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

const (
	// Context key to store Clerk claims/session information
	ClerkSessionContextKey = "clerk_session"
	ClerkUserIDContextKey  = "clerk_user_id"
)

// ContextKey defines a type for context keys to avoid collisions.
type ContextKey string

// UserIDKey is the key used to store the Clerk User ID in the Gin context.
const UserIDKey ContextKey = "clerkUserID"

// ClerkMiddleware creates Gin middleware for Clerk authentication.
func ClerkMiddleware(clerkClient clerk.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the session token from the Authorization header
		sessionToken := c.GetHeader("Authorization")
		if sessionToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// The token is usually prefixed with "Bearer ", remove it
		if len(sessionToken) > 7 && sessionToken[:7] == "Bearer " {
			sessionToken = sessionToken[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Verify the session token
		sessionClaims, err := clerkClient.VerifyToken(sessionToken)
		if err != nil {
			log.Printf("Error verifying Clerk session token: %v", err)
			// Differentiate between invalid token and other errors if needed
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			return
		}

		// Check if session is active (optional but recommended)
		// if sessionClaims.Expiry.Before(time.Now()) { // This check might be handled by VerifyToken already
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
		// 	return
		// }

		// Set the user ID in the context for downstream handlers
		c.Set(string(UserIDKey), sessionClaims.Subject) // Subject usually holds the User ID
		log.Printf("Clerk Auth successful for user: %s", sessionClaims.Subject)

		// Continue to the next handler
		c.Next()
	}
}

// GetClerkUserID retrieves the Clerk User ID from the Gin context.
// It returns the UserID and true if found, otherwise an empty string and false.
func GetClerkUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		log.Println("Warning: Clerk User ID not found in context")
		return "", false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		log.Println("Warning: Clerk User ID in context is not a string")
		return "", false
	}

	return userIDStr, true
}