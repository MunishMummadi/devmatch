package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"gin/internal/models" // Ensure this matches your module path

	_ "github.com/mattn/go-sqlite3"
)

// DBService encapsulates database operations.
type DBService struct {
	DB *sql.DB
}

// NewDBService creates a new DBService.
func NewDBService(db *sql.DB) *DBService {
	return &DBService{DB: db}
}

// ConnectDB establishes a connection to the SQLite database and ensures the schema exists.
func ConnectDB(databasePath string) (*sql.DB, error) {
	// Note: While databasePath comes from config (trusted), directly concatenating
	// into DSN isn't ideal. Consider validating the path format rigorously.
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&cache=shared&mode=rwc", databasePath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Printf("Error opening database '%s': %v", databasePath, err) // Removed newline
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Ping the database to verify connection
	pingCtx, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	err = db.PingContext(pingCtx)
	if err != nil {
		db.Close()                                                       // Close the connection if ping fails
		log.Printf("Error pinging database '%s': %v", databasePath, err) // Removed newline
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Printf("Database connection established successfully to %s.", databasePath) // Removed newline

	// Ensure the necessary tables exist
	schemaCtx, cancelSchema := context.WithTimeout(context.Background(), 15*time.Second) // Longer timeout for schema ops
	defer cancelSchema()
	err = ensureSchema(schemaCtx, db)
	if err != nil {
		db.Close()                                            // Close the connection if schema setup fails
		log.Printf("Error ensuring database schema: %v", err) // Removed newline
		return nil, fmt.Errorf("failed to ensure database schema: %w", err)
	}

	return db, nil
}

// ensureSchema creates the database tables if they don't exist.
// Note: This executes multiple DDL statements. If one fails mid-way,
// previous ones are not automatically rolled back by this function.
// For robust migrations, consider a dedicated library (e.g., migrate, sql-migrate).
func ensureSchema(ctx context.Context, db *sql.DB) error {
	schemaSQL := `
	PRAGMA foreign_keys = ON; -- Ensure foreign keys are enforced

	-- Users Table --
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		clerk_user_id TEXT UNIQUE NOT NULL,
		username TEXT,
		picture_url TEXT,
		bio TEXT,
		github_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_users_clerk_user_id ON users(clerk_user_id);

	CREATE TRIGGER IF NOT EXISTS trigger_users_update_updated_at
	AFTER UPDATE ON users FOR EACH ROW
	WHEN OLD.updated_at = NEW.updated_at -- Avoid infinite loop if updated_at is explicitly set
	BEGIN
		UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;

	-- Swipes Table --
	CREATE TABLE IF NOT EXISTS swipes (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		swiper_user_id TEXT NOT NULL,
		swiped_user_id TEXT NOT NULL,
		direction TEXT NOT NULL CHECK(direction IN ('like', 'dislike')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (swiper_user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (swiped_user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (swiper_user_id, swiped_user_id)
	);
	CREATE INDEX IF NOT EXISTS idx_swipes_swiper_id ON swipes(swiper_user_id);
	CREATE INDEX IF NOT EXISTS idx_swipes_swiped_id ON swipes(swiped_user_id);

	-- Conversations Table --
	CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL -- Tracks last activity/message
	);

	-- Conversation Participants Join Table --
	CREATE TABLE IF NOT EXISTS conversation_participants (
		conversation_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		PRIMARY KEY (conversation_id, user_id),
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_conv_participants_user_id ON conversation_participants(user_id);

	-- Messages Table --
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		conversation_id TEXT NOT NULL,
		sender_user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (sender_user_id) REFERENCES users(id) ON DELETE CASCADE -- Assuming sender must exist
	);
	CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
	CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages(sent_at);

	-- Trigger to update conversation updated_at on new message --
	CREATE TRIGGER IF NOT EXISTS trigger_update_conversation_on_message
	AFTER INSERT ON messages FOR EACH ROW
	BEGIN
		UPDATE conversations SET updated_at = NEW.sent_at WHERE id = NEW.conversation_id;
	END;
	`

	log.Println("Ensuring database schema...")
	_, err := db.ExecContext(ctx, schemaSQL)
	if err != nil {
		log.Printf("Error executing schema SQL: %v", err) // Removed newline
		return fmt.Errorf("failed to execute schema SQL: %w", err)
	}
	log.Println("Database schema check complete.")
	return nil
}

// --- User Profile Operations ---

// GetUserProfileByClerkID retrieves a user profile using their Clerk ID.
// Returns sql.ErrNoRows if the user is not found.
func (s *DBService) GetUserProfileByClerkID(ctx context.Context, clerkUserID string) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at
		FROM users
		WHERE clerk_user_id = ?`

	var user models.User
	err := s.DB.QueryRowContext(ctx, query, clerkUserID).Scan(
		&user.ID,
		&user.ClerkUserID,
		&user.Username,
		&user.PictureURL,
		&user.Bio,
		&user.GitHubURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Log the specific case but return the original error for clarity upstream
			log.Printf("No user found with Clerk ID: %q", clerkUserID) // Removed newline
			return nil, err                                            // Return sql.ErrNoRows
		}
		log.Printf("Error querying user by Clerk ID %q: %v", clerkUserID, err) // Removed newline
		return nil, fmt.Errorf("querying user by Clerk ID %q failed: %w", clerkUserID, err)
	}
	return &user, nil
}

// CreateOrUpdateUserProfile creates a new user or updates an existing one based on Clerk User ID.
// Uses a transaction and returns the created or updated user profile.
func (s *DBService) CreateOrUpdateUserProfile(ctx context.Context, user models.User) (*models.User, error) {
	if user.ClerkUserID == "" {
		return nil, errors.New("ClerkUserID is required to create or update profile")
	}

	// Start a transaction
	tx, err := s.DB.BeginTx(ctx, nil) // Use default transaction options
	if err != nil {
		log.Printf("Error starting transaction for ClerkID %q: %v", user.ClerkUserID, err) // Removed newline
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer rollback in case of errors - it's a no-op if Commit() succeeds
	defer tx.Rollback()

	// Use COALESCE to handle nil pointers gracefully in the update part
	upsertQuery := `
		INSERT INTO users (clerk_user_id, username, picture_url, bio, github_url)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (clerk_user_id) DO UPDATE SET
			username = COALESCE(excluded.username, users.username),
			picture_url = COALESCE(excluded.picture_url, users.picture_url),
			bio = COALESCE(excluded.bio, users.bio),
			github_url = COALESCE(excluded.github_url, users.github_url),
			updated_at = CURRENT_TIMESTAMP` // Let trigger handle updated_at if possible, but set here for INSERT case

	_, err = tx.ExecContext(ctx, upsertQuery,
		user.ClerkUserID,
		user.Username,
		user.PictureURL,
		user.Bio,
		user.GitHubURL,
	)
	if err != nil {
		log.Printf("Error executing upsert for ClerkID %q: %v", user.ClerkUserID, err) // Removed newline
		// Rollback is deferred
		return nil, fmt.Errorf("upsert failed for ClerkID %q: %w", user.ClerkUserID, err)
	}

	// Select the user data after upsert
	selectQuery := `
		SELECT id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at
		FROM users
		WHERE clerk_user_id = ?`

	var createdOrUpdatedUser models.User
	err = tx.QueryRowContext(ctx, selectQuery, user.ClerkUserID).Scan(
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
		// This shouldn't happen if the upsert succeeded, but handle it defensively
		log.Printf("Error selecting user after upsert for ClerkID %q: %v", user.ClerkUserID, err) // Removed newline
		// Rollback is deferred
		return nil, fmt.Errorf("failed to select user after upsert for ClerkID %q: %w", user.ClerkUserID, err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction for ClerkID %q: %v", user.ClerkUserID, err) // Removed newline
		// Rollback is deferred, but commit failed
		return nil, fmt.Errorf("failed to commit transaction for ClerkID %q: %w", user.ClerkUserID, err)
	}

	return &createdOrUpdatedUser, nil
}

// TODO: Implement other DB operations (Get Swipes, Create Swipe, Get Conversations, Create Message, etc.)
// Remember to use context and handle potential sql.ErrNoRows appropriately.
