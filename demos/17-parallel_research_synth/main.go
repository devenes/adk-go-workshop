// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package demonstrates parallel agent execution and synthesis.
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type CityInput struct {
	City string `json:"city"`
}

type ResearchOutput struct {
	Results string `json:"results"`
}

func researchMuseums(ctx tool.Context, input CityInput) (ResearchOutput, error) {
	return ResearchOutput{Results: "Top museums in " + input.City + " include art, history, and science museums."}, nil
}

func researchEvents(ctx tool.Context, input CityInput) (ResearchOutput, error) {
	return ResearchOutput{Results: input.City + " hosts festivals, concerts, and cultural events year-round."}, nil
}

func researchFood(ctx tool.Context, input CityInput) (ResearchOutput, error) {
	return ResearchOutput{Results: input.City + " has diverse cuisine with local specialties and restaurants."}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	museumTool, err := functiontool.New(functiontool.Config{
		Name:        "research_museums",
		Description: "Research museums in a city",
	}, researchMuseums)
	if err != nil {
		log.Fatalf("Failed to create museum tool: %v", err)
	}

	eventsTool, err := functiontool.New(functiontool.Config{
		Name:        "research_events",
		Description: "Research events in a city",
	}, researchEvents)
	if err != nil {
		log.Fatalf("Failed to create events tool: %v", err)
	}

	foodTool, err := functiontool.New(functiontool.Config{
		Name:        "research_food",
		Description: "Research food and cuisine in a city",
	}, researchFood)
	if err != nil {
		log.Fatalf("Failed to create food tool: %v", err)
	}

	museumAgent, err := llmagent.New(llmagent.Config{
		Name:        "museum_agent",
		Model:       model,
		Description: "Research museums.",
		Instruction: "Research and describe museums in the specified city.",
		Tools: []tool.Tool{
			museumTool,
		},
		OutputKey: "museum_result",
	})
	if err != nil {
		log.Fatalf("Failed to create museum agent: %v", err)
	}

	eventsAgent, err := llmagent.New(llmagent.Config{
		Name:        "events_agent",
		Model:       model,
		Description: "Research events.",
		Instruction: "Research and describe events in the specified city.",
		Tools: []tool.Tool{
			eventsTool,
		},
		OutputKey: "events_result",
	})
	if err != nil {
		log.Fatalf("Failed to create events agent: %v", err)
	}

	foodAgent, err := llmagent.New(llmagent.Config{
		Name:        "food_agent",
		Model:       model,
		Description: "Research food.",
		Instruction: "Research and describe food and cuisine in the specified city.",
		Tools: []tool.Tool{
			foodTool,
		},
		OutputKey: "food_result",
	})
	if err != nil {
		log.Fatalf("Failed to create food agent: %v", err)
	}

	parallelBlock, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "parallel_city_research",
			Description: "Run research agents in parallel.",
			SubAgents:   []agent.Agent{museumAgent, eventsAgent, foodAgent},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create parallel agent: %v", err)
	}

	synthesisAgent, err := llmagent.New(llmagent.Config{
		Name:        "synthesis_agent",
		Model:       model,
		Description: "Synthesize research results.",
		Instruction: "Synthesize the museum results from {museum_result}, events from {events_result}, and food from {food_result} into a comprehensive city guide.",
	})
	if err != nil {
		log.Fatalf("Failed to create synthesis agent: %v", err)
	}

	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "parallel_then_synthesize",
			Description: "Run parallel research then synthesize results.",
			SubAgents:   []agent.Agent{parallelBlock, synthesisAgent},
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
