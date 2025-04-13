package services

import (
	"context"
	"fmt"
	"log"
	"gin/internal/config"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// ClerkService handles specific interactions with the Clerk API beyond basic auth middleware.
type ClerkService struct {
	client clerk.Client
	cfg    *config.Config
}

// NewClerkService creates a new instance of ClerkService.
func NewClerkService(client clerk.Client, cfg *config.Config) *ClerkService {
	if client == nil {
		log.Println("Warning: Clerk client provided to NewClerkService is nil. Service may not function correctly.")
		// Depending on use case, might return nil or an error, or allow a nil client
	}
	return &ClerkService{
		client: client,
		cfg:    cfg,
	}
}

// GetClerkUserMetadata fetches user metadata from Clerk using the user ID.
// Note: Often, user claims are already available in the request context via middleware.
// This function demonstrates a direct API call if needed.
func (s *ClerkService) GetClerkUserMetadata(ctx context.Context, userID string) (*clerk.User, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Clerk client is not initialized")
	}

	user, err := s.client.Users().Read(userID)
	if err != nil {
		log.Printf("Error fetching Clerk user data for %s: %v", userID, err)
		// Handle specific Clerk errors if necessary, e.g., user not found
		return nil, fmt.Errorf("failed to fetch user data from Clerk: %w", err)
	}

	log.Printf("Successfully fetched Clerk user data for %s", userID)
	return user, nil
}

// Add other Clerk API interaction functions as needed.