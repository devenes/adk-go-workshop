# Getting Started with the ADK Go Workshop

Welcome! This is a complete hands-on workshop for learning the **Google Agent Development Kit (ADK)** in **Go**.

## Before You Start

You need:
- **Go 1.22+** ([install](https://golang.org/doc/install))
- **Google Gemini API key** ([get free key](https://ai.google.dev/gemini-api/docs/api-key))

## 1. Set Up (2 minutes)

```bash
# Clone the repository
git clone <url> && cd google-adk-workshop

# Set your API key (required for all demos)
export GOOGLE_API_KEY="your-key-here"

# Verify setup
go version
go mod tidy
go build ./...
```

## 2. Run Your First Demo (5 minutes)

```bash
# Run the simplest demo in console mode
go run ./demos/01-hello_web console

# Type a question and press Enter
# Example: "What is the capital of France?"

# Exit with Ctrl+C
```

## 3. Try the Web UI (5 minutes)

```bash
# Run demo 05 with web UI
go run ./demos/05-day_trip_search web

# Browser opens at http://localhost:8765
# Chat directly in the web interface
```

## 4. Explore Demos (20+ minutes)

Each demo teaches a different ADK concept. Progress from **Beginner** → **Advanced** → **Expert**:

### Beginner (learn tools)
```bash
go run ./demos/02-calculator_basics console     # Function tools
go run ./demos/03-custom_tools console          # Custom tools
```

### Intermediate (learn workflows)
```bash
go run ./demos/06-session_memory console        # Session state
go run ./demos/07-sequential_pipeline console   # Sequential workflow
go run ./demos/09-live_weather_nws console      # Live API calls
```

### Advanced (learn patterns)
```bash
go run ./demos/12-agent_as_tool_orchestrator console  # Agent-as-tool
go run ./demos/13-structured_output console           # Structured JSON
```

### Expert (learn advanced features)
```bash
go run ./demos/16-loop_plan_refine console      # Iterative refinement
go run ./demos/17-parallel_research_synth console # Parallel execution
```

## 5. Read the Code

Each demo has a well-commented `main.go`. Open any demo to see:

```bash
# Example: look at the calculator demo
cat demos/02-calculator_basics/main.go

# Key sections to understand:
# 1. Creating a model (lines ~30-40)
# 2. Defining tools (lines ~50-70)
# 3. Creating an agent (lines ~72-85)
# 4. Running the launcher (lines ~87-92)
```

## 6. Modify a Demo

Try editing a demo to experiment:

```bash
# Copy demo 01 and customize it
cp -r demos/01-hello_web demos/my-agent
cd demos/my-agent

# Edit main.go - change the instruction or add a tool
nano main.go

# Run your version
go run . console
```

## 7. Run Tests

Verify that all demos are correctly set up:

```bash
# Run all smoke tests
go test ./...

# Test a specific demo
go test ./demos/01-hello_web/...
```

## 8. Check Security

Before shipping code, scan for accidental API key leaks:

```bash
go run scripts/check_api_key_leaks.go
```

## Next Steps

- **Understand the concepts:** See [`CONVERSION_SUMMARY.md`](CONVERSION_SUMMARY.md) for Python ↔ Go patterns
- **Deploy:** Read [`DEPLOY.md`](DEPLOY.md) for Cloud Run and production tips
- **Learn deeply:** Check [`LEARNING_DEEP_DIVE.md`](LEARNING_DEEP_DIVE.md)
- **Integrate:** See [`INTEGRATIONS.md`](INTEGRATIONS.md) for MCP, OpenAPI, etc.
- **Reference docs:** [ADK Documentation](https://google.github.io/adk-docs/)

## Common Issues

### "go: command not found"
You need to install Go. See [golang.org/doc/install](https://golang.org/doc/install)

### "ERROR: No API key"
Set your API key:
```bash
export GOOGLE_API_KEY="your-actual-key-here"
```

### "connection refused"
The demo is waiting for you to type a message. For web UI, check that port 8765 is not in use:
```bash
# On macOS/Linux: find what's using port 8765
lsof -i :8765

# On Windows
netstat -ano | findstr :8765
```

### "build failed"
Make sure you're in the repo root and run:
```bash
go mod tidy
go build ./...
```

## Learning Path (Recommended)

**Day 1 - Basics (2 hours):**
- `01-hello_web` — Minimal agent
- `02-calculator_basics` — Function tools
- `03-custom_tools` — Custom tools

**Day 2 - State & Memory (2 hours):**
- `04-static_kb_rag` — Knowledge base
- `06-session_memory` — Session state
- `09-live_weather_nws` — Live HTTP

**Day 3 - Workflows (2 hours):**
- `07-sequential_pipeline` — Sequential agents
- `08-sequential_state_shared` — State passing
- `11-multi_agent_coordinator` — Multi-agent

**Day 4 - Advanced (2 hours):**
- `12-agent_as_tool_orchestrator` — Agent-as-tool
- `13-structured_output` — JSON schemas
- `15-structured_persona_research` — Combined patterns

**Day 5 - Expert (2 hours):**
- `16-loop_plan_refine` — Iterative refinement
- `17-parallel_research_synth` — Parallel execution

## Cheat Sheet

```bash
# Run a demo
go run ./demos/XX-name [console|web]

# Build all demos
go build ./...

# Test all demos
go test ./...

# Check code for API key leaks
go run scripts/check_api_key_leaks.go

# Show demo code
cat demos/XX-name/main.go

# Run in background
go run ./demos/XX-name console &

# Build binary
go build -o my-agent ./demos/XX-name
./my-agent console
```

## Get Help

- 📖 Read the code: `demos/XX-name/main.go`
- 🧪 Check tests: `demos/XX-name/main_test.go`
- 📚 ADK docs: [google.github.io/adk-docs](https://google.github.io/adk-docs/)
- 🐛 Bug? See [`PRESENTER_GUIDE.md`](PRESENTER_GUIDE.md) for fallbacks

---

**Ready to build agents in Go?** Start with `demos/01-hello_web` and have fun! 🚀
