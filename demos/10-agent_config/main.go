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

// Package demonstrates agent configuration with multiple tools (Go equivalent of YAML config demo).
package main

import (
	"context"
	"log"
	"math/rand"
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

type DieInput struct {
	Sides int `json:"sides"`
}

type DieOutput struct {
	Result int `json:"result"`
}

type NumbersInput struct {
	Numbers []int `json:"numbers"`
}

type PrimeCheckOutput struct {
	Primes []int `json:"primes"`
}

func rollDie(ctx tool.Context, input DieInput) (DieOutput, error) {
	result := rand.Intn(input.Sides) + 1
	return DieOutput{Result: result}, nil
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func checkPrime(ctx tool.Context, input NumbersInput) (PrimeCheckOutput, error) {
	var primes []int
	for _, n := range input.Numbers {
		if isPrime(n) {
			primes = append(primes, n)
		}
	}
	return PrimeCheckOutput{Primes: primes}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	dieTool, err := functiontool.New(functiontool.Config{
		Name:        "roll_die",
		Description: "Roll a die with a specified number of sides",
	}, rollDie)
	if err != nil {
		log.Fatalf("Failed to create die tool: %v", err)
	}

	primeTool, err := functiontool.New(functiontool.Config{
		Name:        "check_prime",
		Description: "Check which numbers in a list are prime",
	}, checkPrime)
	if err != nil {
		log.Fatalf("Failed to create prime tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "dice_prime_agent",
		Model:       model,
		Description: "Agent that can roll dice and check for prime numbers.",
		Instruction: "You are a helpful assistant with tools for rolling dice and checking if numbers are prime. Use these tools to help the user.",
		Tools: []tool.Tool{
			dieTool,
			primeTool,
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
