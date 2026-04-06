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

// Package demonstrates the agent-as-a-tool pattern.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/adk/tool/geminitool"
)

type EmptyInput struct{}

type YearOutput struct {
	Year int `json:"year"`
}

func getCurrentYear(ctx tool.Context, _ EmptyInput) (YearOutput, error) {
	return YearOutput{Year: time.Now().Year()}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	searchSpecialist, err := llmagent.New(llmagent.Config{
		Name:        "search_specialist",
		Model:       model,
		Description: "Specialist agent for research using Google Search.",
		Instruction: "You are a research specialist. Use Google Search to find information about the topic.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create search specialist: %v", err)
	}

	yearTool, err := functiontool.New(functiontool.Config{
		Name:        "get_current_year",
		Description: "Get the current year",
	}, getCurrentYear)
	if err != nil {
		log.Fatalf("Failed to create year tool: %v", err)
	}

	searchTool := agenttool.New(searchSpecialist, &agenttool.Config{
		SkipSummarization: true,
	})

	orchestrator, err := llmagent.New(llmagent.Config{
		Name:        "research_orchestrator",
		Model:       model,
		Description: "Orchestrator that uses search and utility tools.",
		Instruction: "You are a research orchestrator. Use the search tool for research and the year tool for temporal context.",
		Tools: []tool.Tool{
			yearTool,
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
