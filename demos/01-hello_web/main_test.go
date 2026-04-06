package main

import (
	"context"
	"os"
	"testing"

	"google.golang.org/genai"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
)

func TestAgentCreation(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: GOOGLE_API_KEY not set")
	}

	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "hello_workshop_agent",
		Model:       model,
		Description: "A minimal agent that greets the user.",
		Instruction: "You are a friendly assistant that greets the user warmly.",
	})
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	if a.Name() != "hello_workshop_agent" {
		t.Errorf("Expected agent name 'hello_workshop_agent', got %s", a.Name())
	}
}
