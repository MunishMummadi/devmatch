package services

import (
	"context"
	"log"
	"net/http"

	"github.com/MunishMummadi/devmatch/server/internal/config"

	"github.com/google/go-github/v59/github"
)

// GitHubService handles interactions with the GitHub API.
type GitHubService struct {
	client *github.Client
	cfg    *config.Config
}

// NewGitHubService creates a new instance of GitHubService.
// It initializes the GitHub client, potentially using an OAuth token if available in config.
func NewGitHubService(cfg *config.Config) *GitHubService {
	var tc *http.Client
	// Example: Using a personal access token from config
	// Uncomment and adjust if you have a token in your config
	/*
		if cfg.GitHubToken != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: cfg.GitHubToken},
			)
			tc = oauth2.NewClient(context.Background(), ts)
			log.Println("GitHub client initialized with authentication.")
		} else {
			log.Println("GitHub client initialized without authentication (rate limits apply).")
		}
	*/
	// For now, initialize without authentication
	log.Println("GitHub client initialized without authentication (rate limits apply).")
	client := github.NewClient(tc)

	return &GitHubService{
		client: client,
		cfg:    cfg,
	}
}

// GetUserData fetches basic user data from GitHub for the given username.
func (s *GitHubService) GetUserData(ctx context.Context, username string) (*github.User, error) {
	user, _, err := s.client.Users.Get(ctx, username)
	if err != nil {
		log.Printf("Error fetching GitHub user data for %s: %v", username, err)
		return nil, err
	}
	log.Printf("Successfully fetched GitHub user data for %s", username)
	return user, nil
}

// Add more methods here as needed, e.g., GetUserRepos, GetRepoDetails, etc.