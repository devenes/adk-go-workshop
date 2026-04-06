# Documentation Index

Welcome to the Google ADK Go Workshop! Here's a guide to all available documentation.

## 🚀 Getting Started

**New to this workshop?**

- Start here: **[`GETTING_STARTED.md`](GETTING_STARTED.md)** — 5-minute setup and first demo
- Quick overview: **[`README.md`](README.md)** — Project overview and demo table
- Migration guide: **[`PYTHON_TO_GO_MIGRATION.md`](PYTHON_TO_GO_MIGRATION.md)** — For Python developers

## 📚 Learning

### Structured Learning Paths

- **[`CURRICULUM.md`](CURRICULUM.md)** — Beginner to expert learning progression
- **[`COURSE_BEGINNER_TO_EXPERT.md`](COURSE_BEGINNER_TO_EXPERT.md)** — Unified day-long course
- **[`LEARNING_DEEP_DIVE.md`](LEARNING_DEEP_DIVE.md)** — Advanced topics and patterns
- **[`ARCHITECTURE.md`](ARCHITECTURE.md)** — System architecture and Mermaid diagrams

### Reference Materials

- **[`CONVERSION_SUMMARY.md`](CONVERSION_SUMMARY.md)** — Complete Python-to-Go conversion details
- **[`demos/CHECKPOINTS.md`](demos/CHECKPOINTS.md)** — Demo verification checklist
- **[`INTEGRATIONS.md`](INTEGRATIONS.md)** — MCP, OpenAPI, and external integrations
- **[`RUBRIC.md`](RUBRIC.md)** — Learner self-evaluation checklist

## 🎯 Practical Guides

### Running Demos

```bash
# Run any demo (console or web mode)
go run ./demos/01-hello_web console
go run ./demos/05-day_trip_search web

# See all 17 demos with descriptions in README.md
```

### Testing & Deployment

- **[`EVAL.md`](EVAL.md)** — Evaluation sets and schema validation
- **[`DEPLOY.md`](DEPLOY.md)** — Cloud Run deployment and production setup
- **[`PRESENTER_GUIDE.md`](PRESENTER_GUIDE.md)** — For instructors and workshop leaders

## 📖 Code Organization

### Demos (17 progressive examples)

All demos follow the same structure:

```
demos/
├── 01-hello_web/              # Beginner: minimal agent
│   ├── main.go               # Agent code (runnable)
│   └── main_test.go          # Smoke tests
├── 02-calculator_basics/      # Beginner: function tools
├── ...
└── 17-parallel_research_synth/ # Expert: parallel execution
    ├── main.go
    └── main_test.go
```

Each `main.go` can be run with:

```bash
go run ./demos/XX-name console  # Terminal UI
go run ./demos/XX-name web      # Browser UI
```

### Supporting Code

- **`scripts/check_api_key_leaks.go`** — Security scanner for CI/CD
- **`gemini-streamlit-cloudrun/main.go`** — Web server (alternative to demos)
- **`adk-go/`** — Local copy of ADK library (vendored)

## 🔧 Configuration & Build

- **`go.mod`** — Go module definition
- **`go.sum`** — Dependency lock file
- **`.github/workflows/ci.yml`** — GitHub Actions CI/CD pipeline
- **`gemini-streamlit-cloudrun/Dockerfile`** — Docker container build

## 🗺️ Learning Roadmap

### Day 1: Foundations (4 hours)

- Read: `GETTING_STARTED.md`
- Run: Demos 01, 02, 03
- Concept: Basic agents and tools
- File: `demos/02-calculator_basics/main.go`

### Day 2: State & Memory (4 hours)

- Run: Demos 04, 05, 06, 09
- Concept: Session state, knowledge bases, live APIs
- File: `demos/06-session_memory/main.go`

### Day 3: Workflows (4 hours)

- Run: Demos 07, 08, 10, 11
- Concept: Sequential and multi-agent patterns
- File: `demos/07-sequential_pipeline/main.go`

### Day 4: Advanced Patterns (4 hours)

- Run: Demos 12, 13, 14, 15
- Concept: Structured output, agent-as-tool, HITL
- File: `demos/13-structured_output/main.go`

### Day 5: Expert Topics (4 hours)

- Run: Demos 16, 17
- Read: `LEARNING_DEEP_DIVE.md`
- Concept: Iterative refinement, parallel execution
- File: `demos/16-loop_plan_refine/main.go`

## 🎓 Common Questions

**Q: I'm new to Go, should I learn Go first?**
A: Not necessary. Start with `GETTING_STARTED.md`, and you'll learn Go alongside ADK patterns.

**Q: How do I modify a demo?**
A: Copy the demo to a new folder, edit `main.go`, and run it. See `GETTING_STARTED.md` for an example.

**Q: Can I deploy this to production?**
A: Yes! See [`DEPLOY.md`](DEPLOY.md) for Cloud Run and Kubernetes setup.

**Q: How do I use this with my own tools?**
A: Look at `demos/02-calculator_basics/main.go` for function tools, or `demos/12-agent_as_tool_orchestrator/main.go` for sub-agents.

**Q: What if I know Python and want to understand the migration?**
A: Read [`PYTHON_TO_GO_MIGRATION.md`](PYTHON_TO_GO_MIGRATION.md) for side-by-side comparisons.

## 📋 Quick Command Reference

```bash
# Setup
go version                          # Check Go installation
go mod tidy                         # Fetch dependencies
go build ./...                      # Compile all demos

# Running
go run ./demos/XX-name console      # Terminal chat interface
go run ./demos/XX-name web          # Browser UI (opens port 8765)

# Testing
go test ./...                       # Run all tests
go test ./demos/01-hello_web/...    # Test specific demo

# Verification
go run scripts/check_api_key_leaks.go  # Security check

# Building
go build -o my-agent ./demos/XX-name   # Create standalone binary
docker build -t my-app ./gemini-streamlit-cloudrun  # Docker image
```

## 📞 Get Help

| Resource                                       | Use for                                |
| ---------------------------------------------- | -------------------------------------- |
| `GETTING_STARTED.md`                           | First-time setup and running demos     |
| `README.md`                                    | Project overview and demo descriptions |
| `demos/XX-name/main.go`                        | Understanding specific patterns        |
| `PYTHON_TO_GO_MIGRATION.md`                    | Converting Python knowledge to Go      |
| `LEARNING_DEEP_DIVE.md`                        | Advanced topics and best practices     |
| `DEPLOY.md`                                    | Production deployment                  |
| [ADK Docs](https://google.github.io/adk-docs/) | Official ADK documentation             |
| `PRESENTER_GUIDE.md`                           | For instructors                        |

## 🔗 External Links

- **[Google ADK Documentation](https://google.github.io/adk-docs/)** — Official docs
- **[Google GenAI Go SDK](https://pkg.go.dev/google.golang.org/genai)** — Gemini API reference
- **[ADK Crash Course](https://codelabs.developers.google.com/onramp/)** — Free Google Codelabs
- **[Go Documentation](https://golang.org/doc/)** — Go language reference

## 📝 Document Overview

| Document                    | Audience          | Focus               | Time    |
| --------------------------- | ----------------- | ------------------- | ------- |
| `GETTING_STARTED.md`        | Everyone          | Setup and first run | 10 min  |
| `README.md`                 | Everyone          | Overview and demos  | 5 min   |
| `PYTHON_TO_GO_MIGRATION.md` | Python devs       | Pattern translation | 20 min  |
| `CONVERSION_SUMMARY.md`     | Maintainers       | Migration details   | 30 min  |
| `CURRICULUM.md`             | Learners          | Structured paths    | 5 min   |
| `LEARNING_DEEP_DIVE.md`     | Advanced learners | Deep topics         | 60+ min |
| `DEPLOY.md`                 | DevOps/Engineers  | Production setup    | 30 min  |
| `PRESENTER_GUIDE.md`        | Instructors       | Teaching tips       | 20 min  |
| `ARCHITECTURE.md`           | Architects        | System design       | 15 min  |

## ✅ Verification Checklist

After reading this documentation:

- [ ] I can run `go run ./demos/01-hello_web console`
- [ ] I understand the difference between demos (beginner vs expert)
- [ ] I can read and understand `demos/02-calculator_basics/main.go`
- [ ] I know where to find help (this index!)
- [ ] I'm ready to explore more demos or modify one

---

**Ready to learn the ADK?** Start with [`GETTING_STARTED.md`](GETTING_STARTED.md) 🚀
