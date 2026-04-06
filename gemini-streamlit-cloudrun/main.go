// Package main implements a simple web app using the ADK with HTTP API.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

var agentLoader agent.Loader

func init() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	webAgent, err := llmagent.New(llmagent.Config{
		Name:        "web_assistant",
		Model:       model,
		Description: "General purpose web assistant",
		Instruction: "You are a helpful assistant answering questions.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agentLoader = agent.NewSingleLoader(webAgent)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/health", handleHealth)

	// Also expose the web UI launcher API
	// Users can interact via the standard ADK web interface or via simple HTTP endpoints below
	mux.HandleFunc("/api/chat", handleChat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Gemini web app on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
	<title>ADK Gemini Web App</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; max-width: 900px; margin: 0 auto; padding: 20px; }
		h1 { color: #1f2937; }
		.container { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
		.card { border: 1px solid #e5e7eb; border-radius: 8px; padding: 20px; }
		textarea { width: 100%; height: 120px; padding: 10px; border: 1px solid #d1d5db; border-radius: 4px; font-family: monospace; }
		button { padding: 10px 20px; background: #3b82f6; color: white; border: none; cursor: pointer; border-radius: 4px; margin-top: 10px; }
		button:hover { background: #2563eb; }
		.output { margin-top: 15px; padding: 15px; background: #f9fafb; border-radius: 4px; white-space: pre-wrap; word-wrap: break-word; max-height: 300px; overflow-y: auto; }
		.loading { color: #6b7280; font-style: italic; }
		.error { color: #dc2626; }
		.info { color: #4b5563; font-size: 14px; line-height: 1.6; }
	</style>
</head>
<body>
	<h1>🤖 ADK Gemini Web Application</h1>
	<p class="info">Powered by Google Gemini 2.5 Flash and the Google ADK</p>

	<div class="container">
		<div class="card">
			<h2>Ask a Question</h2>
			<textarea id="query" placeholder="Ask me anything..."></textarea>
			<button onclick="askQuestion()">Ask Gemini</button>
			<div id="result" style="display:none;"></div>
		</div>

		<div class="card">
			<h2>About</h2>
			<p class="info">
				This application demonstrates the Google Agent Development Kit (ADK) with Gemini LLM.
			</p>
			<p class="info">
				<strong>Features:</strong>
			</p>
			<ul class="info">
				<li>Powered by Gemini 2.5 Flash</li>
				<li>Built with Google ADK</li>
				<li>Supports Google Search grounding</li>
				<li>Simple HTTP API</li>
			</ul>
		</div>
	</div>

	<script>
		async function askQuestion() {
			const query = document.getElementById('query').value.trim();
			const resultDiv = document.getElementById('result');

			if (!query) {
				alert('Please enter a question');
				return;
			}

			resultDiv.innerHTML = '<p class="loading">⏳ Gemini is thinking...</p>';
			resultDiv.style.display = 'block';

			try {
				const response = await fetch('/api/chat', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ query: query })
				});

				const data = await response.json();
				if (data.error) {
					resultDiv.innerHTML = '<p class="error">❌ ' + escapeHtml(data.error) + '</p>';
				} else {
					resultDiv.innerHTML = '<p>✅ ' + escapeHtml(data.response) + '</p>';
				}
			} catch (err) {
				resultDiv.innerHTML = '<p class="error">❌ Error: ' + escapeHtml(err.message) + '</p>';
			}
		}

		function escapeHtml(text) {
			const div = document.createElement('div');
			div.textContent = text;
			return div.innerHTML;
		}

		// Allow Ctrl+Enter to submit
		document.getElementById('query').addEventListener('keypress', (e) => {
			if (e.key === 'Enter' && e.ctrlKey) {
				askQuestion();
			}
		});
	</script>
</body>
</html>
	`)
}

type ChatRequest struct {
	Query string `json:"query"`
}

type ChatResponse struct {
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, ChatResponse{Error: "Invalid request format"}, http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		respondJSON(w, ChatResponse{Error: "Query cannot be empty"}, http.StatusBadRequest)
		return
	}

	// NOTE: For a real implementation, you would use runner.New() to execute the agent
	// and capture its response. This is a simplified placeholder that shows the API structure.
	// To fully integrate, implement the launcher's HTTP API or use runner.Run()

	response := "The ADK web server is running. To use the full agent functionality with LLM calls, run the demos with the ADK launcher."

	respondJSON(w, ChatResponse{Response: response}, http.StatusOK)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "ok", "version": "1.0"}, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
