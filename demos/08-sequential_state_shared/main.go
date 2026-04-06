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

// Package demonstrates sequential agents with state sharing via OutputKey and placeholders.
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

	pickSpotAgent, err := llmagent.New(llmagent.Config{
		Name:        "pick_destination",
		Model:       model,
		Description: "Picks a destination from the user's request.",
		Instruction: "Extract and return ONLY the destination city/place name from the user's request. Be concise.",
		OutputKey:   "destination",
	})
	if err != nil {
		log.Fatalf("Failed to create pick spot agent: %v", err)
	}

	navigateAgent, err := llmagent.New(llmagent.Config{
		Name:        "navigate_there",
		Model:       model,
		Description: "Provides navigation details to a destination.",
		Instruction: "The user wants to navigate to {destination}. Provide helpful directions and travel tips for this location.",
	})
	if err != nil {
		log.Fatalf("Failed to create navigate agent: %v", err)
	}

	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "sequential_state_shared",
			Description: "Sequential pipeline with state sharing between agents.",
			SubAgents:   []agent.Agent{pickSpotAgent, navigateAgent},
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
