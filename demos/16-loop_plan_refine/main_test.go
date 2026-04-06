package main

import (
	"context"
	"os"
	"testing"

	"google.golang.org/genai"

	"google.golang.org/adk/model/gemini"
)

func TestAgent(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: GOOGLE_API_KEY not set")
	}

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	if model == nil {
		t.Fatal("Model is nil")
	}
}
