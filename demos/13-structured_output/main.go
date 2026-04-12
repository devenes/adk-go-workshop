// Package demonstrates structured JSON output using output schema.
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
)

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	outputSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name": {
				Type:        genai.TypeString,
				Description: "The name of the city",
			},
			"country": {
				Type:        genai.TypeString,
				Description: "The country where the city is located",
			},
			"population_band": {
				Type:        genai.TypeString,
				Description: "The population range (e.g., '1-5 million', '5-10 million')",
			},
			"highlights": {
				Type:        genai.TypeArray,
				Description: "Key highlights and attractions of the city",
				Items: &genai.Schema{
					Type: genai.TypeString,
				},
			},
		},
		Required: []string{"name", "country", "population_band", "highlights"},
	}

	a, err := llmagent.New(llmagent.Config{
		Name:         "structured_city_agent",
		Model:        model,
		Description:  "Agent that provides structured city information.",
		Instruction:  "You are a city information assistant. When asked about a city, respond with structured data including name, country, population band, and key highlights.",
		OutputSchema: outputSchema,
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
