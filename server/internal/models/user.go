package models

import (
	"time"
)

// User represents the user profile stored in the database.
type User struct {
	ID          string    `json:"id" db:"id"`                       // Assuming UUID or string ID from DB
	ClerkUserID string    `json:"-" db:"clerk_user_id"`             // Usually not exposed directly in API responses
	Username    *string   `json:"username,omitempty" db:"username"` // Use pointers for optional fields
	PictureURL  *string   `json:"pictureUrl,omitempty" db:"picture_url"`
	Bio         *string   `json:"bio,omitempty" db:"bio"`
	GitHubURL   *string   `json:"githubUrl,omitempty" db:"github_url"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateUserProfileRequest defines the expected payload for creating/updating a user profile.
// Often similar to the User model but might exclude server-set fields like ID, timestamps.
type CreateUserProfileRequest struct {
	Username   *string `json:"username"`
	PictureURL *string `json:"pictureUrl"`
	Bio        *string `json:"bio"`
	GitHubURL  *string `json:"githubUrl"`
	// ClerkUserID will be added from the authenticated session, not the request body
}
