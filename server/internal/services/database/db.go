package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MunishMummadi/devmatch/server/internal/models" // Use new module path
	_ "github.com/mattn/go-sqlite3"                           // SQLite driver
)

// DBService encapsulates database operations.
type DBService struct {
	DB *sql.DB
}

// NewDBService creates a new DBService.
func NewDBService(db *sql.DB) *DBService {
	return &DBService{DB: db}
}

// ConnectDB establishes a connection to the SQLite database.
func ConnectDB(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return nil, err
	}

	// Configure connection pool (optional but recommended)
	// db.SetMaxOpenConns(10) // Example: Set max open connections
	// db.SetMaxIdleConns(5)  // Example: Set max idle connections
	// db.SetConnMaxLifetime(time.Hour) // Example: Set connection max lifetime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		log.Printf("Error pinging database: %v\n", err)
		return nil, err
	}

	log.Printf("Database connection established successfully to %s.\n", databasePath)
	return db, nil
}

// --- User Profile Operations ---

// GetUserProfileByClerkID retrieves a user profile using their Clerk ID.
func (s *DBService) GetUserProfileByClerkID(ctx context.Context, clerkUserID string) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at
		FROM users
		WHERE clerk_user_id = ?`

	var user models.User
	var username, pictureURL, bio, githubURL sql.NullString

	err := s.DB.QueryRowContext(ctx, query, clerkUserID).Scan(
		&user.ID,
		&user.ClerkUserID,
		&username,
		&pictureURL,
		&bio,
		&githubURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No user found with Clerk ID: %s", clerkUserID)
			return nil, err
		}
		log.Printf("Error querying user by Clerk ID: %v", err)
		return nil, err
	}

	user.Username = sqlNullStringToStringPtr(username)
	user.PictureURL = sqlNullStringToStringPtr(pictureURL)
	user.Bio = sqlNullStringToStringPtr(bio)
	user.GitHubURL = sqlNullStringToStringPtr(githubURL)

	return &user, nil
}

// CreateOrUpdateUserProfile creates a new user or updates an existing one based on Clerk User ID.
func (s *DBService) CreateOrUpdateUserProfile(ctx context.Context, user models.User) (*models.User, error) {
	if user.ClerkUserID == "" {
		return nil, fmt.Errorf("ClerkUserID is required to create or update profile")
	}

	query := `
		INSERT INTO users (clerk_user_id, username, picture_url, bio, github_url)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (clerk_user_id) DO UPDATE SET
			username = excluded.username,
			picture_url = excluded.picture_url,
			bio = excluded.bio,
			github_url = excluded.github_url,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at`

	var createdOrUpdatedUser models.User
	err := s.DB.QueryRowContext(ctx, query,
		user.ClerkUserID,
		user.Username,
		user.PictureURL,
		user.Bio,
		user.GitHubURL,
	).Scan(
		&createdOrUpdatedUser.ID,
		&createdOrUpdatedUser.ClerkUserID,
		&createdOrUpdatedUser.Username,
		&createdOrUpdatedUser.PictureURL,
		&createdOrUpdatedUser.Bio,
		&createdOrUpdatedUser.GitHubURL,
		&createdOrUpdatedUser.CreatedAt,
		&createdOrUpdatedUser.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error in CreateOrUpdateUserProfile: %v\n", err)
		return nil, err
	}

	return &createdOrUpdatedUser, nil
}

// GetUserProfileByID retrieves a user profile using their internal DB ID.
func (s *DBService) GetUserProfileByID(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at
		FROM users
		WHERE id = ?`

	var user models.User
	var username, pictureURL, bio, githubURL sql.NullString

	err := s.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.ClerkUserID,
		&username,
		&pictureURL,
		&bio,
		&githubURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No user found with ID: %s", userID)
			return nil, err
		}
		log.Printf("Error querying user by ID: %v", err)
		return nil, err
	}

	user.Username = sqlNullStringToStringPtr(username)
	user.PictureURL = sqlNullStringToStringPtr(pictureURL)
	user.Bio = sqlNullStringToStringPtr(bio)
	user.GitHubURL = sqlNullStringToStringPtr(githubURL)

	return &user, nil
}

// EditUserProfile updates specific fields of a user's profile.
// Note: This assumes the caller has already verified authorization.
func (s *DBService) EditUserProfile(ctx context.Context, userID string, updates models.CreateUserProfileRequest) (*models.User, error) {
	// Build the update query dynamically based on provided fields
	// IMPORTANT: Only update fields that are not nil in the 'updates' request
	query := "UPDATE users SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argPlaceholders := []string{}

	if updates.Username != nil {
		args = append(args, *updates.Username)
		argPlaceholders = append(argPlaceholders, fmt.Sprintf("username = ?%d", len(args)))
	}
	if updates.PictureURL != nil {
		args = append(args, *updates.PictureURL)
		argPlaceholders = append(argPlaceholders, fmt.Sprintf("picture_url = ?%d", len(args)))
	}
	if updates.Bio != nil {
		args = append(args, *updates.Bio)
		argPlaceholders = append(argPlaceholders, fmt.Sprintf("bio = ?%d", len(args)))
	}
	if updates.GitHubURL != nil {
		args = append(args, *updates.GitHubURL)
		argPlaceholders = append(argPlaceholders, fmt.Sprintf("github_url = ?%d", len(args)))
	}

	// Ensure at least one field was provided for update
	if len(args) == 0 {
		log.Printf("EditUserProfile called with no fields to update for user %s", userID)
		// Return the existing profile without making DB changes
		return s.GetUserProfileByID(ctx, userID)
		// Or return an error: return nil, fmt.Errorf("no fields provided for update")
	}

	query += ", " + strings.Join(argPlaceholders, ", ")
	args = append(args, userID)
	query += fmt.Sprintf(" WHERE id = ?%d RETURNING id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at", len(args))

	var updatedUser models.User
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(
		&updatedUser.ID,
		&updatedUser.ClerkUserID,
		&updatedUser.Username,
		&updatedUser.PictureURL,
		&updatedUser.Bio,
		&updatedUser.GitHubURL,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error updating user: User with ID %s not found", userID)
			return nil, err
		}
		log.Printf("Error updating user profile for ID %s: %v", userID, err)
		return nil, err
	}

	return &updatedUser, nil
}

// GetRandomUsers fetches a list of random users, excluding the specified user ID and those already swiped by the current user.
func (s *DBService) GetRandomUsers(ctx context.Context, currentUserID string, limit int) ([]models.User, error) {
	// Updated query to exclude self and already swiped users
	query := `
		SELECT u.id, u.clerk_user_id, u.username, u.picture_url, u.bio, u.github_url, u.created_at, u.updated_at
		FROM users u
		WHERE u.id != $1
		AND NOT EXISTS (
			SELECT 1
			FROM swipes s
			WHERE s.swiper_id = $1 AND s.swiped_id = u.id
		)
		ORDER BY RANDOM() -- Consider alternatives for performance on very large tables
		LIMIT $2`

	rows, err := s.DB.QueryContext(ctx, query, currentUserID, limit)
	if err != nil {
		log.Printf("Error querying random users (excluding %s and swiped): %v", currentUserID, err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var bio, githubURL, pictureURL sql.NullString
		err := rows.Scan(
			&user.ID,
			&user.ClerkUserID,
			&user.Username,
			&pictureURL,
			&bio,
			&githubURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning random user row: %v", err)
			return nil, err
		}
		user.PictureURL = sqlNullStringToStringPtr(pictureURL)
		user.Bio = sqlNullStringToStringPtr(bio)
		user.GitHubURL = sqlNullStringToStringPtr(githubURL)

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating random user rows: %v", err)
		return nil, err
	}

	if len(users) == 0 {
		log.Printf("No suitable random users found (excluding %s and swiped)", currentUserID)
		// Return empty slice, not an error
	}

	return users, nil
}

// SaveSwipe records a swipe action in the database, ignoring duplicates.
func (s *DBService) SaveSwipe(ctx context.Context, swipe models.Swipe) error {
	// Use ON CONFLICT to gracefully handle duplicate swipes
	query := `
		INSERT INTO swipes (swiper_id, swiped_id, direction)
		VALUES ($1, $2, $3)
		ON CONFLICT (swiper_id, swiped_id) DO NOTHING` // Ignore if the swipe already exists

	result, err := s.DB.ExecContext(ctx, query, swipe.SwiperID, swipe.SwipedID, swipe.Direction)
	if err != nil {
		// Log error even if it's handled by ON CONFLICT, in case it's a different error
		log.Printf("Error saving swipe from %s to %s (%s): %v", swipe.SwiperID, swipe.SwipedID, swipe.Direction, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Log if checking rows affected fails, but don't necessarily return error
		log.Printf("Could not determine rows affected for swipe from %s to %s: %v", swipe.SwiperID, swipe.SwipedID, err)
	} else if rowsAffected > 0 {
		log.Printf("Swipe saved: %s -> %s (%s)", swipe.SwiperID, swipe.SwipedID, swipe.Direction)
	} else {
		log.Printf("Swipe ignored (duplicate?): %s -> %s (%s)", swipe.SwiperID, swipe.SwipedID, swipe.Direction)
	}

	return nil
}

// CheckForMatch checks if a reciprocal 'like' exists after a 'like' swipe.
// Returns true if a match is found, false otherwise.
func (s *DBService) CheckForMatch(ctx context.Context, swiperID, swipedID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM swipes
			WHERE swiper_id = ? AND swiped_id = ? AND direction = ?
		)`

	var exists bool
	// We are checking if the 'swiped' user has previously 'liked' the current 'swiper'
	err := s.DB.QueryRowContext(ctx, query, swipedID, swiperID, models.SwipeLike).Scan(&exists)
	if err != nil {
		log.Printf("Error checking for match between %s and %s: %v", swiperID, swipedID, err)
		return false, err
	}

	log.Printf("Match check: %s -> %s. Reciprocal like exists: %t", swiperID, swipedID, exists)
	return exists, nil
}

// --- Chat Implementation ---

// GetConversationsForUser retrieves conversations a user is part of.
func (s *DBService) GetConversationsForUser(ctx context.Context, userID string) ([]models.Conversation, error) {
	query := `
		SELECT c.id, c.created_at, c.updated_at
		FROM conversations c
		JOIN participants p ON c.id = p.conversation_id
		WHERE p.user_id = $1
		ORDER BY c.updated_at DESC` // Or created_at, depending on desired order

	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying conversations for user %s: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		if err := rows.Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			log.Printf("Error scanning conversation row for user %s: %v", userID, err)
			// Decide whether to return partial results or fail entirely
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating conversation rows for user %s: %v", userID, err)
		return nil, err
	}

	log.Printf("Retrieved %d conversations for user %s", len(conversations), userID)
	return conversations, nil
}

// GetMessagesForConversation retrieves messages for a specific conversation.
// Assumes authorization (user is part of conversation) is checked beforehand.
func (s *DBService) GetMessagesForConversation(ctx context.Context, conversationID string, limit, offset int) ([]models.Message, error) {
	// Added limit and offset for pagination
	query := `
		SELECT id, conversation_id, sender_id, content, timestamp
		FROM messages
		WHERE conversation_id = $1
		ORDER BY timestamp DESC -- Show newest messages first
		LIMIT $2 OFFSET $3`

	rows, err := s.DB.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		log.Printf("Error querying messages for conversation %s: %v", conversationID, err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.Timestamp); err != nil {
			log.Printf("Error scanning message row for conversation %s: %v", conversationID, err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating message rows for conversation %s: %v", conversationID, err)
		return nil, err
	}

	log.Printf("Retrieved %d messages for conversation %s (limit %d, offset %d)", len(messages), conversationID, limit, offset)
	// Reverse the slice so the oldest message is first in the returned slice (more typical for display)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

// SaveMessage saves a new chat message and updates the conversation timestamp.
// Assumes authorization is checked beforehand.
func (s *DBService) SaveMessage(ctx context.Context, message *models.Message) error {
	// Use a transaction to ensure atomicity
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction for SaveMessage: %v", err)
		return err
	}
	defer tx.Rollback() // Rollback if any step fails

	// 1. Insert the message
	// Ensure timestamp and ID are set before saving
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}
	// Assuming message.ID is generated elsewhere (e.g., UUID) before calling this
	if message.ID == "" {
		// If not generated elsewhere, you might need to generate it here or use DB default
		log.Println("Warning: SaveMessage called with empty message ID")
		// For now, let's return an error or generate one if needed
		return errors.New("message ID cannot be empty") // Or generate UUID here
	}

	msgQuery := `
		INSERT INTO messages (id, conversation_id, sender_id, content, timestamp)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, msgQuery, message.ID, message.ConversationID, message.SenderID, message.Content, message.Timestamp)
	if err != nil {
		log.Printf("Error inserting message into conversation %s: %v", message.ConversationID, err)
		return err // Rollback will happen via defer
	}

	// 2. Update the conversation's updated_at timestamp
	convUpdateQuery := `UPDATE conversations SET updated_at = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, convUpdateQuery, message.Timestamp, message.ConversationID)
	if err != nil {
		// Log the error, but maybe the message insert is more critical?
		// Depending on requirements, you might decide to commit the message even if this fails.
		log.Printf("Error updating conversation %s timestamp after new message: %v (message was inserted)", message.ConversationID, err)
		// For now, let's rollback if this fails too to keep things consistent.
		return err // Rollback will happen via defer
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction for SaveMessage (conversation %s): %v", message.ConversationID, err)
		return err
	}

	log.Printf("Message %s saved successfully in conversation %s", message.ID, message.ConversationID)
	return nil
}

// CreateConversation creates a new conversation and adds participants.
func (s *DBService) CreateConversation(ctx context.Context, participantIDs []string) (*models.Conversation, error) {
	if len(participantIDs) < 2 {
		return nil, errors.New("conversation requires at least two participants")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction for CreateConversation: %v", err)
		return nil, err
	}
	defer tx.Rollback() // Rollback if anything fails

	// 1. Insert into conversations table
	now := time.Now()
	var conv models.Conversation
	convQuery := `INSERT INTO conversations (created_at, updated_at) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, convQuery, now, now).Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting new conversation: %v", err)
		return nil, err
	}

	// 2. Insert participants into the participants table
	partQuery := `INSERT INTO participants (conversation_id, user_id) VALUES ($1, $2)`
	stmt, err := tx.PrepareContext(ctx, partQuery)
	if err != nil {
		log.Printf("Error preparing participant insert statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	for _, userID := range participantIDs {
		if _, err := stmt.ExecContext(ctx, conv.ID, userID); err != nil {
			log.Printf("Error inserting participant %s into conversation %s: %v", userID, conv.ID, err)
			// Check for specific errors like duplicate participant if needed
			return nil, err // Rollback
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction for CreateConversation: %v", err)
		return nil, err
	}

	log.Printf("Conversation %s created successfully with participants: %v", conv.ID, participantIDs)
	return &conv, nil
}

// IsUserParticipant checks if a user is a participant in a given conversation.
func (s *DBService) IsUserParticipant(ctx context.Context, userID, conversationID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM participants WHERE conversation_id = $1 AND user_id = $2)`
	var exists bool
	err := s.DB.QueryRowContext(ctx, query, conversationID, userID).Scan(&exists)
	if err != nil {
		// Don't log ErrNoRows as an error here, it just means false.
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error checking participation for user %s in conversation %s: %v", userID, conversationID, err)
		}
		// Return the error only if it's unexpected.
		return false, err
	}
	return exists, nil
}

// --- Chat Operations ---

// GetConversationsByUserID retrieves all conversations a user is part of.
// Includes basic info like conversation ID, participants, and last update time.
func (s *DBService) GetConversationsByUserID(ctx context.Context, userID string) ([]models.Conversation, error) {
	// This query retrieves conversation IDs and timestamps for conversations
	// involving the given userID. It then needs a second step (or a more complex join/aggregation)
	// to get all participants for each conversation.
	query := `
		SELECT c.id, c.created_at, c.updated_at
		FROM conversations c
		JOIN participants p ON c.id = p.conversation_id
		WHERE p.user_id = ?
		ORDER BY c.updated_at DESC` // Order by most recently active

	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying conversations for user %s: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning conversation row for user %s: %v", userID, err)
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating conversation rows for user %s: %v", userID, err)
		return nil, err
	}

	if len(conversations) == 0 {
		log.Printf("No conversations found for user %s", userID)
		return []models.Conversation{}, nil // Return empty slice, not error
	}

	// Create a map for easy lookup
	conversationsMap := make(map[string]*models.Conversation)
	var conversationIDs []string
	for i := range conversations {
		conv := &conversations[i] // Get pointer to the conversation in the slice
		conversationsMap[conv.ID] = conv
		conversationIDs = append(conversationIDs, conv.ID)
	}

	// Now fetch participants for each conversation
	participantQuery := `SELECT user_id FROM participants WHERE conversation_id = ?`
	for _, convID := range conversationIDs { // Iterate in the original order (most recent first)
		partRows, err := s.DB.QueryContext(ctx, participantQuery, convID)
		if err != nil {
			log.Printf("Error fetching participants for conversation %s: %v", convID, err)
			// Decide how to handle partial failure - skip this conversation? return error?
			// For now, we log and continue, the conversation will lack participant IDs.
			continue
		}

		var participantIDs []string
		for partRows.Next() {
			var pUserID string // Use a different variable name to avoid shadowing
			if err := partRows.Scan(&pUserID); err != nil {
				log.Printf("Error scanning participant user ID for conversation %s: %v", convID, err)
				continue // Skip this participant
			}
			participantIDs = append(participantIDs, pUserID)
		}
		// It's crucial to check for errors *after* the loop
		if err = partRows.Err(); err != nil {
			log.Printf("Error iterating participants for conversation %s: %v", convID, err)
			// Continue to next conversation even if participant iteration had errors for this one
		}
		partRows.Close() // Close rows inside the loop

		if conv, ok := conversationsMap[convID]; ok {
			conv.ParticipantIDs = participantIDs
		}
	}

	return conversations, nil
}

// CreateMessage saves a new chat message and updates the conversation's timestamp.
func (s *DBService) CreateMessage(ctx context.Context, message models.Message) (*models.Message, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction for CreateMessage: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	// 1. Insert the message
	// Assume SentAt should be Timestamp based on current Message struct
	newMessage := message
	newMessage.Timestamp = time.Now()

	// In a real scenario, you might want more robust ID generation
	newMessage.ID = "some-generated-id" // Replace with actual ID generation logic

	msgQuery := `
		INSERT INTO messages (id, conversation_id, sender_id, content, timestamp)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, conversation_id, sender_id, content, timestamp`
	err = tx.QueryRowContext(ctx, msgQuery, newMessage.ID, newMessage.ConversationID, newMessage.SenderID, newMessage.Content, newMessage.Timestamp).Scan(
		&newMessage.ID,
		&newMessage.ConversationID,
		&newMessage.SenderID,
		&newMessage.Content,
		&newMessage.Timestamp,
	)
	if err != nil {
		log.Printf("Error inserting message into conversation %s by user %s: %v", newMessage.ConversationID, newMessage.SenderID, err)
		return nil, err
	}

	// 2. Update the conversation's updated_at timestamp
	convUpdateQuery := `UPDATE conversations SET updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = tx.ExecContext(ctx, convUpdateQuery, newMessage.ConversationID)
	if err != nil {
		// Log the error but potentially allow message creation to succeed anyway?
		// Or rollback? Rolling back for consistency.
		log.Printf("Error updating conversation %s timestamp after new message: %v", newMessage.ConversationID, err)
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction for CreateMessage: %v", err)
		return nil, err
	}

	log.Printf("Message %s saved in conversation %s by user %s", newMessage.ID, newMessage.ConversationID, newMessage.SenderID)
	return &newMessage, nil
}

// GetMessagesByConversationID retrieves messages for a specific conversation, ordered by time.
// Supports basic pagination using limit and offset.
func (s *DBService) GetMessagesByConversationID(ctx context.Context, conversationID string, limit, offset int) ([]models.Message, error) {
	// Ensure limit is reasonable
	if limit <= 0 || limit > 100 { // Set a max limit
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, conversation_id, sender_id, content, timestamp
		FROM messages
		WHERE conversation_id = ?
		ORDER BY timestamp ASC -- Fetch newest first commonly
		LIMIT ? OFFSET ?`

	rows, err := s.DB.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		log.Printf("Error querying messages for conversation %s: %v", conversationID, err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.Content,
			&msg.Timestamp,
		)
		if err != nil {
			log.Printf("Error scanning message row for conversation %s: %v", conversationID, err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating message rows for conversation %s: %v", conversationID, err)
		return nil, err
	}

	// Reverse the slice because we queried ASC for pagination but usually want chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	log.Printf("Fetched %d messages for conversation %s (limit %d, offset %d)", len(messages), conversationID, limit, offset)
	return messages, nil
}

// --- Favorite Operations ---

// IsFavorite checks if favoriteUserID is in userID's favorites list.
func (s *DBService) IsFavorite(ctx context.Context, userID, favoriteUserID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = ? AND favorite_user_id = ?)`
	var exists bool
	err := s.DB.QueryRowContext(ctx, query, userID, favoriteUserID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user %s is favorite of user %s: %v", favoriteUserID, userID, err)
		return false, err
	}
	return exists, nil
}

// AddFavorite adds a user to another user's favorites list.
// It uses INSERT OR IGNORE to avoid errors if the favorite already exists.
func (s *DBService) AddFavorite(ctx context.Context, userID, favoriteUserID string) error {
	query := `INSERT OR IGNORE INTO favorites (user_id, favorite_user_id) VALUES (?, ?)`
	_, err := s.DB.ExecContext(ctx, query, userID, favoriteUserID)
	if err != nil {
		log.Printf("Error adding favorite: User %s -> Fav %s: %v", userID, favoriteUserID, err)
		return err
	}
	log.Printf("Added favorite: User %s -> Fav %s", userID, favoriteUserID)
	return nil
}

// RemoveFavorite removes a user from another user's favorites list.
func (s *DBService) RemoveFavorite(ctx context.Context, userID, favoriteUserID string) error {
	query := `DELETE FROM favorites WHERE user_id = ? AND favorite_user_id = ?`
	result, err := s.DB.ExecContext(ctx, query, userID, favoriteUserID)
	if err != nil {
		log.Printf("Error removing favorite: User %s -> Fav %s: %v", userID, favoriteUserID, err)
		return err
	}
	rowsAffected, _ := result.RowsAffected() // Check if a row was actually deleted
	log.Printf("Removed favorite: User %s -> Fav %s (Rows affected: %d)", userID, favoriteUserID, rowsAffected)
	return nil
}

// GetFavoritesByUserID retrieves the list of users favorited by a specific user.
func (s *DBService) GetFavoritesByUserID(ctx context.Context, userID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.clerk_user_id, u.username, u.picture_url, u.bio, u.github_url, u.created_at, u.updated_at
		FROM users u
		JOIN favorites f ON u.id = f.favorite_user_id
		WHERE f.user_id = ?
		ORDER BY f.created_at DESC` // Order by when they were favorited

	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying favorites for user %s: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var favoriteUsers []models.User
	for rows.Next() {
		var user models.User
		var username, pictureURL, bio, githubURL sql.NullString
		err := rows.Scan(
			&user.ID,
			&user.ClerkUserID,
			&username,
			&pictureURL,
			&bio,
			&githubURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning favorite user row for user %s: %v", userID, err)
			return nil, err
		}
		user.Username = sqlNullStringToStringPtr(username)
		user.PictureURL = sqlNullStringToStringPtr(pictureURL)
		user.Bio = sqlNullStringToStringPtr(bio)
		user.GitHubURL = sqlNullStringToStringPtr(githubURL)
		favoriteUsers = append(favoriteUsers, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating favorite user rows for user %s: %v", userID, err)
		return nil, err
	}

	if len(favoriteUsers) == 0 {
		log.Printf("No favorites found for user %s", userID)
		// Return empty slice, not error
	}

	return favoriteUsers, nil
}

// --- Chat Implementation: Fetching Combined Data for API ---

// GetChatDataForUser retrieves the complete chat data (summaries and initial messages)
// for a user, formatted for the frontend API response.
func (s *DBService) GetChatDataForUser(ctx context.Context, userID string, messagesLimitPerConv int) (*models.ChatDataResponse, error) {
	resp := &models.ChatDataResponse{
		Conversations: []models.FrontendConversationSummary{},
		Messages:      make(map[string][]models.FrontendMessage),
	}

	// 1. Get conversations involving the user, along with the other participant's details
	convQuery := `
        SELECT
            c.id AS conversation_id,
            other_p.user_id AS other_user_id,
            other_user.username AS other_username,
            other_user.picture_url AS other_picture_url
            -- Add c.title if you have a title column in conversations
        FROM conversations c
        JOIN participants p ON c.id = p.conversation_id
        JOIN participants other_p ON c.id = other_p.conversation_id AND other_p.user_id != ?
        JOIN users other_user ON other_p.user_id = other_user.id
        WHERE p.user_id = ?
        ORDER BY c.updated_at DESC; -- Order by most recently updated conversation
    `

	rows, err := s.DB.QueryContext(ctx, convQuery, userID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resp, nil // No conversations, return empty response
		}
		log.Printf("Error querying conversations for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	type ConversationInfo struct {
		ConversationID  string
		OtherUserID     string
		OtherUsername   sql.NullString
		OtherPictureURL sql.NullString
		LastMessageText string    // Will be filled later
		LastMessageTime time.Time // Will be filled later
	}

	var convInfos []ConversationInfo
	convIDs := []string{}

	for rows.Next() {
		var info ConversationInfo
		if err := rows.Scan(
			&info.ConversationID,
			&info.OtherUserID,
			&info.OtherUsername,
			&info.OtherPictureURL,
		); err != nil {
			log.Printf("Error scanning conversation row for user %s: %v", userID, err)
			return nil, fmt.Errorf("failed to scan conversation info: %w", err)
		}
		convInfos = append(convInfos, info)
		convIDs = append(convIDs, info.ConversationID)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating conversation rows for user %s: %v", userID, err)
		return nil, fmt.Errorf("error iterating conversation rows: %w", err)
	}

	if len(convIDs) == 0 {
		return resp, nil // No conversations found
	}

	// 2. Get the latest N messages for all these conversations
	msgQueryPlaceholders := make([]string, len(convIDs))
	args := make([]interface{}, len(convIDs)+1)
	for i, id := range convIDs {
		msgQueryPlaceholders[i] = "?"
		args[i] = id
	}
	args[len(convIDs)] = messagesLimitPerConv

	msgQuery := fmt.Sprintf(`
        WITH RankedMessages AS (
            SELECT
                m.id,
                m.conversation_id,
                m.sender_id,
                m.content,
                m.timestamp,
                ROW_NUMBER() OVER(PARTITION BY m.conversation_id ORDER BY m.timestamp DESC) as rn
            FROM messages m
            WHERE m.conversation_id IN (%s)
        )
        SELECT
            rm.id,
            rm.conversation_id,
            rm.sender_id,
            rm.content,
            rm.timestamp
        FROM RankedMessages rm
        WHERE rn <= ?
        ORDER BY rm.conversation_id, rm.timestamp ASC; -- Order messages within conversation chronologically
    `, strings.Join(msgQueryPlaceholders, ","))

	msgRows, err := s.DB.QueryContext(ctx, msgQuery, args...)
	if err != nil {
		log.Printf("Error querying messages for conversations (%v): %v", convIDs, err)
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer msgRows.Close()

	// 3. Process messages and map them to conversations
	rawMessagesByConvID := make(map[string][]models.Message)
	lastMessages := make(map[string]models.Message) // Store the actual last message for summary

	for msgRows.Next() {
		var msg models.Message
		if err := msgRows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.Content,
			&msg.Timestamp,
		); err != nil {
			log.Printf("Error scanning message row: %v", err)
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		rawMessagesByConvID[msg.ConversationID] = append(rawMessagesByConvID[msg.ConversationID], msg)

		// Keep track of the latest message for each conversation
		existingLast, ok := lastMessages[msg.ConversationID]
		if !ok {
			// This is the first message seen for this conversation in this batch
			lastMessages[msg.ConversationID] = msg
		} else {
			// We have seen a message before, check if the current one is newer
			if msg.Timestamp.After(existingLast.Timestamp) {
				lastMessages[msg.ConversationID] = msg
			}
		}
	}
	if err = msgRows.Err(); err != nil {
		log.Printf("Error after iterating message rows: %v", err)
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	// 4. Assemble the final response structure
	convInfoMap := make(map[string]ConversationInfo)
	for _, info := range convInfos {
		convInfoMap[info.ConversationID] = info
	}

	for _, convID := range convIDs { // Iterate in the original order (most recent first)
		info, ok := convInfoMap[convID]
		if !ok {
			continue
		} // Should not happen

		rawMessages := rawMessagesByConvID[convID]
		formattedMessages := make([]models.FrontendMessage, 0, len(rawMessages))
		contactName := "Unknown"
		if info.OtherUsername.Valid {
			contactName = info.OtherUsername.String
		}

		// Format messages for this conversation
		for _, msg := range rawMessages {
			senderName := contactName // Assume sender is the other user
			if msg.SenderID == userID {
				senderName = "You"
			}
			formattedMessages = append(formattedMessages, models.FrontendMessage{
				From:      senderName,
				Text:      msg.Content,
				Timestamp: formatRelativeTime(msg.Timestamp),    // Use the helper
				Type:      models.GetMessageType(msg.SenderID, userID), // Use the helper
			})
		}

		// Populate conversation summary
		lastMsgText := ""
		lastMsgTimestamp := ""
		if lastMsg, hasLast := lastMessages[convID]; hasLast {
			lastMsgText = lastMsg.Content
			// TODO: Use a better relative time formatter for the summary timestamp (e.g., "1h ago")
			lastMsgTimestamp = formatRelativeTime(lastMsg.Timestamp)
		}

		avatar := "https://via.placeholder.com/40" // Default avatar
		if info.OtherPictureURL.Valid {
			avatar = info.OtherPictureURL.String
		}

		resp.Conversations = append(resp.Conversations, models.FrontendConversationSummary{
			ID:          convID,
			UserID:      userID,
			ContactName: contactName,
			LastMessage: lastMsgText,
			Title:       contactName, // Using contact name as title for now
			Timestamp:   lastMsgTimestamp,
			Avatar:      avatar,
		})

		// Add the formatted messages to the map
		if len(formattedMessages) > 0 {
			resp.Messages[convID] = formattedMessages
		}
	}

	return resp, nil
}

// formatRelativeTime converts a time.Time into a human-readable relative string.
// Example implementation, can be made more sophisticated.
func formatRelativeTime(ts time.Time) string {
	now := time.Now()
	duration := now.Sub(ts)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 48*time.Hour {
		return "yesterday"
	} else {
		return ts.Format("Jan 2") // Older than yesterday, show date
	}
}

// --- User Profile Operations ---

// Helper function to handle nullable strings from DB and return a pointer.
func sqlNullStringToStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String // Return pointer to the string
	}
	return nil // Return nil if NULL
}

// GetConversationParticipants retrieves a list of user IDs participating in a specific conversation.
func (s *DBService) GetConversationParticipants(ctx context.Context, conversationID string) ([]string, error) {
	query := `SELECT user_id FROM conversation_participants WHERE conversation_id = ?`

	rows, err := s.DB.QueryContext(ctx, query, conversationID)
	if err != nil {
		log.Printf("Error querying conversation participants for conversation %s: %v", conversationID, err)
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	var participantIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning participant user ID: %v", err)
			return nil, fmt.Errorf("failed to scan participant ID: %w", err)
		}
		participantIDs = append(participantIDs, userID)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating participant rows for conversation %s: %v", conversationID, err)
		return nil, fmt.Errorf("error iterating participant rows: %w", err)
	}

	if len(participantIDs) == 0 {
		// Consider if this is an error or just an empty conversation.
		// For now, return empty slice and no error, but could log a warning.
		log.Printf("Warning: No participants found for conversation %s", conversationID)
	}

	return participantIDs, nil
}
