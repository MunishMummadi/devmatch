package models

import "time"

// Conversation represents a chat conversation between two or more users.
// TODO: Add fields for Title, LastMessage preview etc. if needed for summaries
type Conversation struct {
	ID        string    `json:"id" db:"id"`                 // Unique identifier for the conversation
	CreatedAt time.Time `json:"createdAt" db:"created_at"` // Timestamp when the conversation was created
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // Timestamp when the conversation was last updated
	// TODO: Consider adding Participants []User or ParticipantIDs []string if needed directly on the model
	ParticipantIDs []string `json:"participantIds,omitempty" db:"-"` // List of user IDs participating in the conversation (populated when needed, not a direct DB column on 'conversations')
}
