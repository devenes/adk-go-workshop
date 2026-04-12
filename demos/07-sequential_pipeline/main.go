// Package demonstrates a sequential agent pipeline (outline → expand).
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
)

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	outlineAgent, err := llmagent.New(llmagent.Config{
		Name:        "outline_drafter",
		Model:       model,
		Description: "Creates a high-level outline.",
		Instruction: "Create a detailed 3-point outline for the user's topic. Return only the outline points.",
	})
	if err != nil {
		log.Fatalf("Failed to create outline agent: %v", err)
	}

	expandAgent, err := llmagent.New(llmagent.Config{
		Name:        "detail_expander",
		Model:       model,
		Description: "Expands an outline into full paragraphs.",
		Instruction: "Expand each point from the outline into a full paragraph. Be detailed and comprehensive.",
	})
	if err != nil {
		log.Fatalf("Failed to create expand agent: %v", err)
	}

	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "sequential_write_pipeline",
			Description: "A sequential pipeline that outlines then expands content.",
			SubAgents:   []agent.Agent{outlineAgent, expandAgent},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create sequential agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(sequentialAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
