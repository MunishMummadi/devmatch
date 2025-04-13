package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseDBURL  string
	ClerkSecretKey string
	GeminiAPIKey   string // Add if needed now
	GinMode        string
	Port           string
}

func LoadConfig() *Config {
	// Attempt to load .env file. If it doesn't exist, environment variables are used.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, loading from environment variables")
	}

	cfg := &Config{
		SupabaseDBURL:  getEnv("SUPABASE_DB_URL", ""),
		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", ""),
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""), // Optional for now
		GinMode:        getEnv("GIN_MODE", "debug"),
		Port:           getEnv("PORT", "8080"), // Default port
	}

	if cfg.SupabaseDBURL == "" {
		log.Fatal("FATAL: SUPABASE_DB_URL environment variable is required")
	}
	if cfg.ClerkSecretKey == "" {
		log.Fatal("FATAL: CLERK_SECRET_KEY environment variable is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Warning: Environment variable %s not set, using fallback '%s'\n", key, fallback)
	return fallback
}
