package models

import "time"

// FrontendConversationSummary represents the summary of a conversation for the frontend list.
type FrontendConversationSummary struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"` // The ID of the user making the request
	ContactName string `json:"contactName"`
	LastMessage string `json:"lastMessage"`
	Title       string `json:"title"`       // Consider how this title is generated or stored
	Timestamp   string `json:"timestamp"` // Formatted timestamp (e.g., "1h ago")
	Avatar      string `json:"avatar"`      // URL to the contact's avatar
}

// FrontendMessage represents a single message formatted for the frontend chat view.
type FrontendMessage struct {
	From      string `json:"from"`      // Name of the sender ("You" or contact name)
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"` // Formatted timestamp (e.g., "10:17 am")
	Type      string `json:"type"`      // "incoming" or "outgoing"
}

// ChatDataResponse is the top-level structure for the chat data returned to the frontend.
type ChatDataResponse struct {
	Conversations []FrontendConversationSummary `json:"conversations"`
	Messages      map[string][]FrontendMessage  `json:"messages"` // Key is ConversationID
}

// Helper function (example) - This might live elsewhere (e.g., utils package)
// You'll need a more robust implementation for relative time formatting.
func FormatRelativeTime(t time.Time) string {
	// Placeholder logic - replace with a proper time formatting library or function
	duration := time.Since(t)
	if duration.Hours() < 24 {
		return t.Format("3:04 pm") // e.g., "10:17 am"
	}
	return t.Format("Jan 2") // e.g., "Apr 12"
	// For relative like "1h ago", you'd need more logic or a library.
}

// Helper function (example) - Determining message type
func GetMessageType(messageSenderID string, currentUserID string) string {
	if messageSenderID == currentUserID {
		return "outgoing"
	}
	return "incoming"
}
