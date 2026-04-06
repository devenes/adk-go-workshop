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

	weatherTool, err := functiontool.New(functiontool.Config{
		Name:        "get_weather",
		Description: "Get weather",
	}, getWeather)
	if err != nil {
		t.Fatalf("Failed to create weather tool: %v", err)
	}

	timeTool, err := functiontool.New(functiontool.Config{
		Name:        "get_current_time",
		Description: "Get time",
	}, getCurrentTime)
	if err != nil {
		t.Fatalf("Failed to create time tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:   "weather_time_workshop_agent",
		Model:  model,
		Tools:  []tool.Tool{weatherTool, timeTool},
	})
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	if a.Name() != "weather_time_workshop_agent" {
		t.Errorf("Expected 'weather_time_workshop_agent', got %s", a.Name())
	}
}
