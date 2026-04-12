// Package demonstrates session-scoped memory using tool state.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type RememberInput struct {
	DisplayName string `json:"display_name"`
}

type RememberOutput struct {
	Message string `json:"message"`
}

type RecallOutput struct {
	DisplayName string `json:"display_name"`
}

func rememberDisplayName(ctx tool.Context, input RememberInput) (RememberOutput, error) {
	ctx.Actions().StateDelta["display_name"] = input.DisplayName
	return RememberOutput{Message: fmt.Sprintf("Remembered your display name: %s", input.DisplayName)}, nil
}

func recallDisplayName(ctx tool.Context, _ interface{}) (RecallOutput, error) {
	val, err := ctx.State().Get("display_name")
	if err != nil {
		return RecallOutput{DisplayName: "Unknown"}, nil
	}
	if name, ok := val.(string); ok {
		return RecallOutput{DisplayName: name}, nil
	}
	return RecallOutput{DisplayName: "Unknown"}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	rememberTool, err := functiontool.New(functiontool.Config{
		Name:        "remember_display_name",
		Description: "Remember the user's display name in session memory",
	}, rememberDisplayName)
	if err != nil {
		log.Fatalf("Failed to create remember tool: %v", err)
	}

	// For the recall tool with no input, we use an empty struct
	type EmptyInput struct{}
	recallTool, err := functiontool.New(functiontool.Config{
		Name:        "recall_display_name",
		Description: "Recall the user's previously remembered display name from session memory",
	}, func(ctx tool.Context, input EmptyInput) (RecallOutput, error) {
		val, err := ctx.State().Get("display_name")
		if err != nil {
			return RecallOutput{DisplayName: "Unknown"}, nil
		}
		if name, ok := val.(string); ok {
			return RecallOutput{DisplayName: name}, nil
		}
		return RecallOutput{DisplayName: "Unknown"}, nil
	})
	if err != nil {
		log.Fatalf("Failed to create recall tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "session_memory_workshop_agent",
		Model:       model,
		Description: "Agent that demonstrates session-scoped memory.",
		Instruction: "You are a personal assistant. You can remember the user's display name and recall it later in the conversation. Ask the user for their display name if you don't have it, and use the tools to store and retrieve it.",
		Tools: []tool.Tool{
			rememberTool,
			recallTool,
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
