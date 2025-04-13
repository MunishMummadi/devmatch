package database

import (
	"context"
	"log"
	"time"

	"gin/internal/models" // Corrected import path

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBService encapsulates database operations.
type DBService struct {
	Pool *pgxpool.Pool
}

// NewDBService creates a new DBService.
func NewDBService(pool *pgxpool.Pool) *DBService {
	return &DBService{Pool: pool}
}

// ConnectDB establishes a connection pool to the database.
func ConnectDB(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Error parsing database URL: %v\n", err)
		return nil, err
	}

	// Optional: Configure pool settings
	config.MaxConns = 10 // Example: Set max connections
	config.MinConns = 2  // Example: Set min connections
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return nil, err
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close() // Close pool if ping fails
		log.Printf("Error pinging database: %v\n", err)
		return nil, err
	}

	log.Println("Database connection pool established successfully.")
	return pool, nil
}

// --- User Profile Operations ---

// GetUserProfileByClerkID retrieves a user profile using their Clerk ID.
func (s *DBService) GetUserProfileByClerkID(ctx context.Context, clerkUserID string) (*models.User, error) {
	query := `
		SELECT id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at
		FROM users
		WHERE clerk_user_id = $1`

	var user models.User
	err := s.Pool.QueryRow(ctx, query, clerkUserID).Scan(
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
		// Use pgx.ErrNoRows to check if user not found specifically
		// Log other errors
		return nil, err
	}
	return &user, nil
}

// CreateOrUpdateUserProfile creates a new user or updates an existing one based on Clerk User ID.
// Assumes the input user model contains the ClerkUserID.
func (s *DBService) CreateOrUpdateUserProfile(ctx context.Context, user models.User) (*models.User, error) {
	// Ensure ClerkUserID is provided
	if user.ClerkUserID == "" {
		return nil, log.Output(1, "ClerkUserID is required to create or update profile") // Simplified error handling
	}

	query := `
		INSERT INTO users (clerk_user_id, username, picture_url, bio, github_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (clerk_user_id) DO UPDATE SET
			username = EXCLUDED.username,
			picture_url = EXCLUDED.picture_url,
			bio = EXCLUDED.bio,
			github_url = EXCLUDED.github_url,
			updated_at = NOW() -- Use NOW() or rely on trigger
		RETURNING id, clerk_user_id, username, picture_url, bio, github_url, created_at, updated_at`

	var createdOrUpdatedUser models.User
	err := s.Pool.QueryRow(ctx, query,
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
		log.Printf("Error in CreateOrUpdateUserProfile: %v\n", err) // Log the error
		return nil, err
	}

	return &createdOrUpdatedUser, nil
}

// TODO: Implement other DB operations (Get Random Users, Conversations, Messages, Swipes, Favorites)
