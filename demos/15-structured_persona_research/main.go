// Package demonstrates combining structured output, AgentTool, and Google Search.
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
)

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	personSpecialist, err := llmagent.New(llmagent.Config{
		Name:        "person_specialist",
		Model:       model,
		Description: "Specialist for researching people.",
		Instruction: "You are a people research specialist. Use Google Search to find information about the person.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create specialist: %v", err)
	}

	outputSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name": {
				Type:        genai.TypeString,
				Description: "Full name of the person",
			},
			"age": {
				Type:        genai.TypeInteger,
				Description: "Estimated age of the person",
			},
			"occupation": {
				Type:        genai.TypeString,
				Description: "Primary occupation or profession",
			},
			"location": {
				Type:        genai.TypeString,
				Description: "Primary location or city of residence",
			},
			"biography": {
				Type:        genai.TypeString,
				Description: "Brief biography or description",
			},
		},
		Required: []string{"name", "occupation", "biography"},
	}

	searchTool := agenttool.New(personSpecialist, &agenttool.Config{
		SkipSummarization: true,
	})

	orchestrator, err := llmagent.New(llmagent.Config{
		Name:         "person_orchestrator",
		Model:        model,
		Description:  "Orchestrator that researches people and returns structured data.",
		Instruction:  "You are a people research orchestrator. Use the search tool to research the person and return structured information about them.",
		OutputSchema: outputSchema,
		Tools: []tool.Tool{
			searchTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(orchestrator),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
