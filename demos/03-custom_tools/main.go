// Package demonstrates custom function tools with weather and time functionality.
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
	"google.golang.org/adk/tool/functiontool"
)

type CityInput struct {
	City string `json:"city"`
}

type WeatherOutput struct {
	Weather string `json:"weather"`
}

type TimeOutput struct {
	Time string `json:"time"`
}

func getWeather(ctx tool.Context, input CityInput) (WeatherOutput, error) {
	// Stub implementation
	if input.City == "New York" {
		return WeatherOutput{Weather: "Partly cloudy, 72°F"}, nil
	}
	return WeatherOutput{Weather: "Weather data unavailable for " + input.City}, nil
}

func getCurrentTime(ctx tool.Context, input CityInput) (TimeOutput, error) {
	loc, err := time.LoadLocation(input.City)
	if err != nil {
		// Fallback to UTC if timezone not found
		loc = time.UTC
	}
	currentTime := time.Now().In(loc).Format("15:04:05 MST")
	return TimeOutput{Time: currentTime}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	weatherTool, err := functiontool.New(functiontool.Config{
		Name:        "get_weather",
		Description: "Get the current weather for a city (stub implementation)",
	}, getWeather)
	if err != nil {
		log.Fatalf("Failed to create weather tool: %v", err)
	}

	timeTool, err := functiontool.New(functiontool.Config{
		Name:        "get_current_time",
		Description: "Get the current time in a specific city timezone",
	}, getCurrentTime)
	if err != nil {
		log.Fatalf("Failed to create time tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "weather_time_workshop_agent",
		Model:       model,
		Description: "Agent that provides weather and time information for cities.",
		Instruction: "You are a weather and time assistant. Use the tools to answer questions about weather and current time in various cities.",
		Tools: []tool.Tool{
			weatherTool,
			timeTool,
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
