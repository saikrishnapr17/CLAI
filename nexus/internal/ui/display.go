package ui

import (
    "fmt"
    "strings"

    "github.com/fatih/color"
)

func ShowBanner() {
    cyan := color.New(color.FgHiCyan)
    cyan.Println("  ███╗   ██╗███████╗██╗  ██╗██╗   ██╗███████╗")
    cyan.Println("  ████╗  ██║██╔════╝╚██╗██╔╝██║   ██║██╔════╝")
    cyan.Println("  ██╔██╗ ██║█████╗   ╚███╔╝ ██║   ██║███████╗")
    cyan.Println("  ██║╚██╗██║██╔══╝   ██╔██╗ ██║   ██║╚════██║")
    cyan.Println("  ██║ ╚████║███████╗██╔╝ ██╗╚██████╔╝███████║")
    cyan.Println("  ╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝")
    cyan.Println("  Natural language terminal · powered by Groq")
}

func ShowCommand(command string) {
    yellow := color.New(color.FgHiYellow)
    // build box
    lines := strings.Split(command, "\n")
    max := 0
    for _, l := range lines {
        if len(l) > max {
            max = len(l)
        }
    }
    border := strings.Repeat("─", max+4)
    yellow.Printf("┌─ Command %s┐\n", border)
    for _, l := range lines {
        yellow.Printf("│  %-*s  │\n", max, l)
    }
    yellow.Printf("└%s┘\n", strings.Repeat("─", max+6))
}

func ShowSuccess(message string) {
    green := color.New(color.FgGreen)
    green.Printf("✓ %s\n", message)
}

func ShowError(message string) {
    red := color.New(color.FgRed)
    red.Printf("✗ %s\n", message)
}

func ShowInfo(message string) {
    cyan := color.New(color.FgHiCyan)
    cyan.Printf("ℹ %s\n", message)
}

func ShowPrompt() {
    cyan := color.New(color.FgHiCyan)
    fmt.Print(cyan.Sprint("nexus ❯ "))
}
