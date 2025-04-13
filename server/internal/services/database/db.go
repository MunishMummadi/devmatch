package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gin/internal/models"
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
			log.Printf("No user found with Clerk ID: %s", clerkUserID)
			return nil, err
		}
		log.Printf("Error querying user by Clerk ID: %v", err)
		return nil, err
	}
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

// TODO: Implement other DB operations (Get Random Users, Conversations, Messages, Swipes, Favorites)
