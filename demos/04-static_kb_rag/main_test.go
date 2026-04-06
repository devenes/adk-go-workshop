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
	model, _ := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{APIKey: apiKey})
	kbTool, _ := functiontool.New(functiontool.Config{Name: "search_company_kb", Description: "Search KB"}, searchCompanyKB)
	a, err := llmagent.New(llmagent.Config{Name: "static_kb_agent", Model: model, Tools: []tool.Tool{kbTool}})
	if err != nil || a.Name() != "static_kb_agent" {
		t.Fatal("Failed to create agent")
	}
}
