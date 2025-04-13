package handlers

import (
	"context"
	"log"
	"net/http"

	"gin/internal/services"

	"github.com/gin-gonic/gin"
)

// GitHubHandler handles API requests related to GitHub data fetching and analysis.
type GitHubHandler struct {
	githubService *services.GitHubService
	geminiService *services.GeminiService 
}

// NewGitHubHandler creates a new GitHubHandler.
func NewGitHubHandler(github *services.GitHubService, gemini *services.GeminiService) *GitHubHandler {
	return &GitHubHandler{
		githubService: github,
		geminiService: gemini, 
	}
}

// GetGitHubData fetches data for a given GitHub username.
// GET /github/:username/data
func (h *GitHubHandler) GetGitHubData(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GitHub username is required"})
		return
	}

	if h.githubService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub service not available"})
		return
	}

	userData, err := h.githubService.GetUserData(context.Background(), username)
	if err != nil {
		log.Printf("Error calling GitHub service for %s: %v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch GitHub data"})
		return
	}

	c.JSON(http.StatusOK, userData) 
}

// SummarizeGitHubData receives data (presumably GitHub profile/repo info) and sends it to Gemini for summarization.
// POST /github/summary
func (h *GitHubHandler) SummarizeGitHubData(c *gin.Context) {
	if h.geminiService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Gemini service is not available"})
		return
	}

	var requestBody map[string]interface{} 
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	textToSummarize := "Placeholder text extracted from request" 

	summary, err := h.geminiService.AnalyzeText(context.Background(), textToSummarize)
	if err != nil {
		log.Printf("Error calling Gemini service for summarization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}