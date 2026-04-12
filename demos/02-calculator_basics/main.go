// Package demonstrates function tools for arithmetic operations.
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
	"google.golang.org/adk/tool/functiontool"
)

type NumbersInput struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type NumbersOutput struct {
	Result float64 `json:"result"`
}

func addNumbers(ctx tool.Context, input NumbersInput) (NumbersOutput, error) {
	return NumbersOutput{Result: input.A + input.B}, nil
}

func multiplyNumbers(ctx tool.Context, input NumbersInput) (NumbersOutput, error) {
	return NumbersOutput{Result: input.A * input.B}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	addTool, err := functiontool.New(functiontool.Config{
		Name:        "add_numbers",
		Description: "Adds two numbers together",
	}, addNumbers)
	if err != nil {
		log.Fatalf("Failed to create add tool: %v", err)
	}

	multiplyTool, err := functiontool.New(functiontool.Config{
		Name:        "multiply_numbers",
		Description: "Multiplies two numbers together",
	}, multiplyNumbers)
	if err != nil {
		log.Fatalf("Failed to create multiply tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "calculator_basics_agent",
		Model:       model,
		Description: "An agent that performs basic arithmetic using tools.",
		Instruction: "You are a calculator assistant. Use the provided tools to add or multiply numbers when asked.",
		Tools: []tool.Tool{
			addTool,
			multiplyTool,
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
