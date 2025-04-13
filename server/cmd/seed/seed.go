package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gin/internal/config" // Assuming config might be needed, adjust if not
	"gin/internal/models"
	"gin/internal/services"
	"gin/internal/services/database"
)

// --- Configuration ---
const dbPath = "devmatch.db"   // Relative path from project root
const numberOfUsersToSeed = 50 // Target number of users

// List of potential GitHub usernames to seed from. Refined list.
var githubUsernames = []string{
	// Prominent Devs (mix)
	"gaearon", "addyosmani", "sindresorhus", "yyx990803", "JakeWharton",
	"evanw", "Rich-Harris", "antfu", "patak-dev", "sebmarkbage",
	"kentcdodds", "getify", "wesbos", "rauchg", "tj",
	// Women in Tech / Design / Data Science
	"sarahdrasner", "cassidoo", "una", "jennybc", "jensimmons",
	"mayuko", "angezanetti", "dianaduj", "noopkat", "sailorhg",
	"sophiebits", "lara_hogan", // Placeholder, might need actual GH username
	"ashedryden", "ericasadun", "tracyhinds",
	// Potentially Younger/Streamer Vibe (Mix of real/plausible patterns)
	"codewithlinda", "techgirlwonder", "anastasiajsdev", "frontendunicorn", "pixelista",
	"coderella", "streamerdev", "webdevjourney", "codingkitty", "datavizdiva",
	"jessfraz", "adriennefriend", "ladyleet",
	"gurlcode", "devchic", "miss_debugger", "syntaxsugar", "girldeveloper",
	// Larger Orgs/Projects (for variety)
	"google", "facebook", "microsoft", "github", "octocat",
	"vercel", "netlify", "docker", "kubernetes", "nodejs",
	// Other popular/diverse profiles
	"torvalds", "paulirish", "fabpot", "unclebob", "LeaVerou",
	"matz", "dhh", "jeresig", "tenderlove", "haacked",
	"shanselman", "btholt", "jaredpalmer", "developit", "egoist",
	"mdo", "fat", "twbs", "hadley", "yihui",
	"maxogden", "feross", "mafintosh", "substack", "chriscoyier",
	"daverupert", "zellwk", "mariusschulz", "elijahmanor", "argyleink",
	// Add ~20 more unique plausible usernames if needed to ensure > 50 unique successful fetches
	"devcommunity", "opensourcefan", "codecademy", "freecodecamp", "theodinproject",
	"hackernoon", "devto", "stackshare", "producthunt", "indiehackers",
	"buildspace", "fastai", "huggingface", "tensorflow", "pytorch",
	"sveltejs", "vuejs", "angular", "reactjs", "nextjs",
}

func main() {
	log.Println("Starting database seeding process...")
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	ctx := context.Background()

	// 1. Connect to Database
	db, err := database.ConnectDB(dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	dbService := database.NewDBService(db)

	// 2. Initialize GitHub Service
	// Passing nil for config as the current NewGitHubService doesn't strictly require it
	// for unauthenticated requests. If you add auth later, you'll need a real config.
	var cfg *config.Config // Or initialize with default/empty if needed
	githubService := services.NewGitHubService(cfg)

	// 3. Fetch and Seed Users
	seededCount := 0
	shuffledUsernames := shuffle(githubUsernames) // Shuffle to get variety

	for _, username := range shuffledUsernames {
		if seededCount >= numberOfUsersToSeed {
			break // Stop once we hit the target number
		}

		log.Printf("Attempting to fetch data for GitHub user: %s", username)

		// Fetch data from GitHub
		githubUser, err := githubService.GetUserData(ctx, username)
		if err != nil {
			log.Printf("⚠️ Failed to fetch data for %s: %v. Skipping.", username, err)
			// Optional: Implement backoff/retry logic here if hitting rate limits
			time.Sleep(1 * time.Second) // Simple delay to potentially avoid rate limits
			continue
		}

		// Map to our User model
		userToSeed := models.User{
			// IMPORTANT: Create a predictable, fake Clerk User ID for seeding purposes
			ClerkUserID: fmt.Sprintf("seed_user_%s", username),
			Username:    githubUser.Login,     // Usually same as username input, but use response value
			PictureURL:  githubUser.AvatarURL, // Pointer fields are handled correctly
			Bio:         githubUser.Bio,
			GitHubURL:   githubUser.HTMLURL,
			// ID, CreatedAt, UpdatedAt will be handled by the DB/service
		}

		// Create or Update user in DB
		createdUser, err := dbService.CreateOrUpdateUserProfile(ctx, userToSeed)
		if err != nil {
			log.Printf("⚠️ Failed to seed user %s (ClerkID: %s): %v", *userToSeed.Username, userToSeed.ClerkUserID, err)
			continue
		}

		log.Printf("✅ Successfully seeded user: %s (DB ID: %s)", *createdUser.Username, createdUser.ID)
		seededCount++
	}

	log.Printf("Database seeding process finished. Seeded %d users.", seededCount)
}

// shuffle randomly shuffles a slice of strings
func shuffle(slice []string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := make([]string, len(slice))
	perm := r.Perm(len(slice))
	for i, randIndex := range perm {
		ret[i] = slice[randIndex]
	}
	return ret
}
