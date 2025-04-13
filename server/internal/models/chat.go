package models

import "time"

// Conversation represents a chat session between users.
type Conversation struct {
	ID        string    `json:"id" db:"id"`                 // Unique identifier for the conversation
	UserIDs   []string  `json:"user_ids" db:"user_ids"`     // Slice of user IDs participating in the chat
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Timestamp when the conversation was created
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Timestamp of the last message or update
	// Add other fields like LastMessageSnippet, UnreadCounts per user, etc. if needed
}

// Message represents a single message within a conversation.
type Message struct {
	ID             string    `json:"id" db:"id"`                           // Unique identifier for the message
	ConversationID string    `json:"conversation_id" db:"conversation_id"` // ID of the conversation this message belongs to
	SenderID       string    `json:"sender_id" db:"sender_id"`             // ID of the user who sent the message
	Content        string    `json:"content" db:"content"`                 // The text content of the message
	SentAt         time.Time `json:"sent_at" db:"sent_at"`                 // Timestamp when the message was sent
	// Add other fields like ReadStatus, MessageType (text, image), etc. if needed
}