// Package demonstrates human-in-the-loop (HITL) for sensitive actions.
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

type EmailInput struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailOutput struct {
	Message string `json:"message"`
}

func sendEmail(ctx tool.Context, input EmailInput) (EmailOutput, error) {
	fmt.Printf("[EMAIL SIMULATION] To: %s\nSubject: %s\nBody: %s\n", input.To, input.Subject, input.Body)
	return EmailOutput{Message: "Email sent (simulated)"}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	emailTool, err := functiontool.New(functiontool.Config{
		Name:        "send_email",
		Description: "Send an email message",
	}, sendEmail)
	if err != nil {
		log.Fatalf("Failed to create email tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "hitl_email_agent",
		Model:       model,
		Description: "Agent that sends emails with human oversight.",
		Instruction: "You are an email assistant. Before using the send_email tool, ALWAYS ask the user for explicit confirmation. Show the email details and wait for the user to approve before sending.",
		Tools: []tool.Tool{
			emailTool,
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
