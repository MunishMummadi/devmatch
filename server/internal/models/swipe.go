package models

import "time"

// SwipeDirection defines the possible swipe actions.
type SwipeDirection string

const (
	SwipeLike    SwipeDirection = "like"    // Represents a positive swipe (e.g., right)
	SwipeDislike SwipeDirection = "dislike" // Represents a negative swipe (e.g., left)
)

// Swipe represents a swipe action performed by one user on another.
type Swipe struct {
	ID        string         `json:"id" db:"id"`                 // Unique identifier for the swipe action
	SwiperID  string         `json:"swiper_id" db:"swiper_id"`   // ID of the user who performed the swipe
	SwipedID  string         `json:"swiped_id" db:"swiped_id"`   // ID of the user who was swiped on
	Direction SwipeDirection `json:"direction" db:"direction"` // The direction of the swipe (like/dislike)
	CreatedAt time.Time      `json:"created_at" db:"created_at"` // Timestamp when the swipe occurred
	// Could add MatchID string if a match is created immediately upon swiping
}