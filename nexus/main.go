package main

import (
    "fmt"
    "os"
    "runtime"
    "strings"

    "github.com/joho/godotenv"

    "github.com/yourusername/nexus/internal/executor"
    "github.com/yourusername/nexus/internal/input"
    "github.com/yourusername/nexus/internal/llm"
    "github.com/yourusername/nexus/internal/ui"
)

func isBlocked(command string) bool {
    lower := strings.ToLower(command)
    blocked := []string{
        "rm -rf /",
        "rm -rf ~",
        "mkfs",
        "dd if=",
        ":(){ :|:& };:",
        "chmod -r 777 /",
    }
    for _, b := range blocked {
        if strings.Contains(lower, b) {
            return true
        }
    }
    return false
}

func main() {
    _ = godotenv.Load()

    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        fmt.Println("Set your Groq API key: export GROQ_API_KEY=your_key")
        fmt.Println("Get a free key at console.groq.com")
        os.Exit(1)
    }

    cwd, _ := os.Getwd()
    home, _ := os.UserHomeDir()
    shell := os.Getenv("SHELL")
    osctx := llm.OSContext{
        OS:         runtime.GOOS,
        Shell:      shell,
        CurrentDir: cwd,
        HomeDir:    home,
    }

    client := llm.NewGroqClient(apiKey)

    ui.ShowBanner()
    ui.ShowInfo("Type your request in plain English. Type 'exit' to quit.")

    for {
        ui.ShowPrompt()
        inputLine, err := input.ReadInput("")
        if err != nil {
            ui.ShowError("Failed to read input: " + err.Error())
            continue
        }
        if inputLine == "" {
            continue
        }
        if inputLine == "exit" || inputLine == "quit" {
            break
        }

        ui.ShowInfo("Thinking...")
        systemPrompt := llm.BuildSystemPrompt(osctx)
        command, err := client.Translate(inputLine, systemPrompt)
        if err != nil {
            if strings.Contains(err.Error(), "timed out") {
                ui.ShowError("Groq API request timed out. Try again later.")
                continue
            }
            ui.ShowError("Failed to reach Groq API: " + err.Error())
            continue
        }

        switch command {
        case "CANNOT_TRANSLATE":
            ui.ShowError("Could not convert to a shell command. Try rephrasing.")
            continue
        case "UNSAFE":
            ui.ShowError("That command was flagged as potentially unsafe.")
            continue
        }

        ui.ShowCommand(command)

        if !input.Confirm("Run this?") {
            ui.ShowInfo("Skipped.")
            continue
        }

        if isBlocked(command) {
            ui.ShowError("Blocked: potentially destructive command.")
            continue
        }

        err = executor.Execute(command)
        if err != nil {
            ui.ShowError("Command failed: " + err.Error())
        } else {
            ui.ShowSuccess("Done.")
        }
    }

    ui.ShowInfo("Goodbye.")
}
