// Package demonstrates live HTTP calls to the National Weather Service API.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

type ForecastOutput struct {
	Forecast string `json:"forecast"`
}

var cityCoords = map[string]struct{ lat, lon string }{
	"New York":    {"40.7128", "-74.0060"},
	"Los Angeles": {"34.0522", "-118.2437"},
	"Chicago":     {"41.8781", "-87.6298"},
	"Houston":     {"29.7604", "-95.3698"},
}

func getLiveForecast(ctx tool.Context, input CityInput) (ForecastOutput, error) {
	city := strings.TrimSpace(input.City)
	coords, ok := cityCoords[city]
	if !ok {
		return ForecastOutput{Forecast: fmt.Sprintf("Forecast data not available for %s. Supported cities: New York, Los Angeles, Chicago, Houston", city)}, nil
	}

	// Get grid point data
	gridURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", coords.lat, coords.lon)
	gridResp, err := http.Get(gridURL)
	if err != nil {
		return ForecastOutput{Forecast: fmt.Sprintf("Error fetching grid data: %v", err)}, nil
	}
	defer gridResp.Body.Close()

	var gridData map[string]interface{}
	if err := json.NewDecoder(gridResp.Body).Decode(&gridData); err != nil {
		return ForecastOutput{Forecast: fmt.Sprintf("Error parsing grid data: %v", err)}, nil
	}

	properties, ok := gridData["properties"].(map[string]interface{})
	if !ok {
		return ForecastOutput{Forecast: "Unable to extract properties from grid data"}, nil
	}

	forecastURL, ok := properties["forecast"].(string)
	if !ok {
		return ForecastOutput{Forecast: "Unable to extract forecast URL"}, nil
	}

	// Get forecast
	forecastResp, err := http.Get(forecastURL)
	if err != nil {
		return ForecastOutput{Forecast: fmt.Sprintf("Error fetching forecast: %v", err)}, nil
	}
	defer forecastResp.Body.Close()

	body, err := io.ReadAll(forecastResp.Body)
	if err != nil {
		return ForecastOutput{Forecast: fmt.Sprintf("Error reading forecast: %v", err)}, nil
	}

	var forecastData map[string]interface{}
	if err := json.Unmarshal(body, &forecastData); err != nil {
		return ForecastOutput{Forecast: fmt.Sprintf("Error parsing forecast: %v", err)}, nil
	}

	if fp, ok := forecastData["properties"].(map[string]interface{}); ok {
		if periods, ok := fp["periods"].([]interface{}); ok && len(periods) > 0 {
			if period, ok := periods[0].(map[string]interface{}); ok {
				shortForecast, _ := period["shortForecast"].(string)
				temperature, _ := period["temperature"].(float64)
				return ForecastOutput{Forecast: fmt.Sprintf("%.0f°F, %s", temperature, shortForecast)}, nil
			}
		}
	}

	return ForecastOutput{Forecast: "Forecast data unavailable"}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	forecastTool, err := functiontool.New(functiontool.Config{
		Name:        "get_live_forecast",
		Description: "Get live weather forecast from the US National Weather Service",
	}, getLiveForecast)
	if err != nil {
		log.Fatalf("Failed to create forecast tool: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "nws_weather_agent",
		Model:       model,
		Description: "Agent that provides live weather forecasts from the National Weather Service.",
		Instruction: "You are a weather forecasting assistant. Use the tool to get live weather forecasts for US cities and provide helpful weather analysis.",
		Tools: []tool.Tool{
			forecastTool,
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
