# Nexus — Natural language terminal

Nexus lets you type plain English and translates it to a shell command (via the Groq API), confirms with you, then executes it.

## Features
- Interactive REPL: prompt `nexus ❯` and type natural language requests.
- Translates requests using Groq (`llama-4-scout-17b-16e-instruct`).
- Confirms before executing and streams command output in real time.
- Safety checks to block destructive commands.

## Requirements
- Go 1.22+

## Setup
1. Copy the example env file and set your key:

```bash
cp .env.example .env
export GROQ_API_KEY=your_groq_api_key_here
```

You can get a key at https://console.groq.com

## Build
```bash
go mod tidy
go build -o nexus .
```

## Run
Start the interactive CLI:

```bash
./nexus
```

Usage notes:
- Type plain English requests, e.g. `list my go files` or `find . -name "*.md"`.
- Type `exit` or `quit` to close.
- If Groq cannot translate a request it will return `CANNOT_TRANSLATE`.

## Safety
Nexus enforces several runtime checks to avoid destructive operations (for example `rm -rf /`, `mkfs`, `dd if=` and common fork-bombs). The CLI will refuse to execute flagged commands and will print a clear error.

## Files of interest
- [main.go](main.go) — REPL loop and wiring
- [internal/llm/groq.go](internal/llm/groq.go) — Groq API client
- [internal/llm/prompt.go](internal/llm/prompt.go) — prompt builder
- [internal/executor/run.go](internal/executor/run.go) — command execution
- [internal/input/reader.go](internal/input/reader.go) — terminal input helpers
- [internal/ui/display.go](internal/ui/display.go) — colored terminal UI

## Next steps / optional features
- Command history (`~/.nexus_history`)
- `--dry-run` and `--explain` flags
- Pipe mode for non-interactive usage

---
If you want, I can add history and `--dry-run` support next.
