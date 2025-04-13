package services

import (
	"context"
	"log"

	"github.com/MunishMummadi/devmatch/server/internal/config"
	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

// ClerkService handles specific interactions with the Clerk API beyond basic auth middleware.
type ClerkService struct {
	cfg *config.Config
}

// NewClerkService creates a new instance of ClerkService.
func NewClerkService(cfg *config.Config) *ClerkService {
	return &ClerkService{
		cfg: cfg,
	}
}

// GetUser fetches user details from Clerk using the user ID.
func (s *ClerkService) GetUser(ctx context.Context, userID string) (*clerk.User, error) {
	clerkUser, err := user.Get(ctx, userID)
	if err != nil {
		log.Printf("Error fetching user from Clerk API for userID %s: %v", userID, err)
		return nil, err
	}
	return clerkUser, nil
}

// Add other Clerk API interaction functions as needed.
