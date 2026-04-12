// Package demonstrates in-memory knowledge base search (static RAG).
package main

import (
	"context"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type SearchInput struct {
	Query string `json:"query"`
}

type SearchOutput struct {
	Results string `json:"results"`
}

var knowledgeBase = map[string]string{
	"billing_cycle":   "Billing cycles occur monthly on the 1st of each month. Invoices are sent via email.",
	"api_rate_limits": "API rate limits are 1000 requests per minute for standard tier, 10000 for premium tier.",
	"data_retention":  "Data is retained for 90 days by default. Extended retention available on request.",
}

func searchCompanyKB(ctx tool.Context, input SearchInput) (SearchOutput, error) {
	query := strings.ToLower(input.Query)
	for key, value := range knowledgeBase {
		if strings.Contains(strings.ToLower(key), query) || strings.Contains(strings.ToLower(value), query) {
			return SearchOutput{Results: value}, nil
		}
	}
	return SearchOutput{Results: "No matching information found in knowledge base."}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	kbTool, err := functiontool.New(functiontool.Config{
		Name:        "search_company_kb",
		Description: "Search the company knowledge base for policies and information",
	}, searchCompanyKB)
	if err != nil {
		log.Fatalf("Failed to create KB tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "static_kb_agent",
		Model:       model,
		Description: "Agent that answers questions by searching a static knowledge base.",
		Instruction: "You are a company support assistant. Use the knowledge base search tool to answer customer questions about policies, billing, and API details.",
		Tools: []tool.Tool{
			kbTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
