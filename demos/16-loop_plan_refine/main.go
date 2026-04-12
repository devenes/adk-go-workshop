// Package demonstrates iterative refinement using LoopAgent.
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type EmptyInput struct{}

type ExitLoopOutput struct {
	Message string `json:"message"`
}

func exitLoop(ctx tool.Context, _ EmptyInput) (ExitLoopOutput, error) {
	ctx.Actions().Escalate = true
	return ExitLoopOutput{Message: "Loop exit triggered"}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	plannerAgent, err := llmagent.New(llmagent.Config{
		Name:        "planner_agent",
		Model:       model,
		Description: "Creates an initial plan.",
		Instruction: "You are a planner. Create a detailed plan with 3 main steps.",
		OutputKey:   "current_plan",
	})
	if err != nil {
		log.Fatalf("Failed to create planner: %v", err)
	}

	criticAgent, err := llmagent.New(llmagent.Config{
		Name:        "critic_agent",
		Model:       model,
		Description: "Critiques the current plan.",
		Instruction: "You are a critic. Analyze the plan at {current_plan} and provide constructive criticism and areas for improvement.",
		OutputKey:   "criticism",
	})
	if err != nil {
		log.Fatalf("Failed to create critic: %v", err)
	}

	exitTool, err := functiontool.New(functiontool.Config{
		Name:        "exit_loop",
		Description: "Exit the refinement loop when plan is good enough",
	}, exitLoop)
	if err != nil {
		log.Fatalf("Failed to create exit tool: %v", err)
	}

	refinerAgent, err := llmagent.New(llmagent.Config{
		Name:        "refiner_agent",
		Model:       model,
		Description: "Refines the plan based on criticism.",
		Instruction: "You are a refiner. Take the plan at {current_plan} and the criticism at {criticism}, then create an improved version. If the plan is good enough, use the exit_loop tool.",
		OutputKey:   "current_plan",
		Tools: []tool.Tool{
			exitTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create refiner: %v", err)
	}

	refinementLoop, err := loopagent.New(loopagent.Config{
		AgentConfig: agent.Config{
			Name:        "refinement_loop",
			Description: "Loop that refines the plan iteratively.",
			SubAgents:   []agent.Agent{criticAgent, refinerAgent},
		},
		MaxIterations: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create loop: %v", err)
	}

	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "iterative_plan_workshop",
			Description: "Sequential pipeline with iterative refinement.",
			SubAgents:   []agent.Agent{plannerAgent, refinementLoop},
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
