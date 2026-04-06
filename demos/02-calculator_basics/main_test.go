package main

import (
	"context"
	"os"
	"testing"

	"google.golang.org/genai"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
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

	addTool, err := functiontool.New(functiontool.Config{
		Name:        "add_numbers",
		Description: "Adds two numbers together",
	}, addNumbers)
	if err != nil {
		t.Fatalf("Failed to create add tool: %v", err)
	}

	multiplyTool, err := functiontool.New(functiontool.Config{
		Name:        "multiply_numbers",
		Description: "Multiplies two numbers together",
	}, multiplyNumbers)
	if err != nil {
		t.Fatalf("Failed to create multiply tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "calculator_basics_agent",
		Model:       model,
		Description: "An agent that performs basic arithmetic using tools.",
		Instruction: "You are a calculator assistant.",
		Tools:       []tool.Tool{addTool, multiplyTool},
	})
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	if a.Name() != "calculator_basics_agent" {
		t.Errorf("Expected agent name 'calculator_basics_agent', got %s", a.Name())
	}
}
