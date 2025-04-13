package models

import (
	"time"
)

// Message represents a single message within a conversation in the database.
type Message struct {
	ID             string    `json:"id" db:"id"`                           // Unique identifier for the message
	ConversationID string    `json:"conversationId" db:"conversation_id"` // ID of the conversation this message belongs to
	SenderID       string    `json:"senderId" db:"sender_id"`             // ID of the user who sent the message
	Content        string    `json:"content" db:"content"`                 // The text content of the message
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`             // Timestamp when the message was sent
}

// FrontendMessage is defined in chat_api.go
// Helper functions FormatRelativeTime and GetMessageType are defined in chat_api.go
