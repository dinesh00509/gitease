package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinesh00509/gitease/internals"
)

const version = "v1.1.0"

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--version" || arg == "-v" {
			fmt.Println("GitEase", version)
			os.Exit(0)
		}
		if arg == "--help" || arg == "-h" {
			fmt.Println("GitEase - Interactive Git Assistant")
			fmt.Println("Usage: gitease [flags]")
			fmt.Println("Navigate and manage Git tasks interactively in your terminal.")
			fmt.Println("Flags:")
			fmt.Println("  --version, -v   Show version")
			fmt.Println("  --help, -h      Show this help message")
			fmt.Println("  --run, -r       Start the interactive CLI")
			os.Exit(0)
		}
		if arg == "--run" || arg == "-r" {
		} else {
			fmt.Println("Unknown flag:", arg)
			fmt.Println("Use --help for usage information")
			os.Exit(1)
		}
	}

	p := tea.NewProgram(internals.InitialModel())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running GitEase:", err)
		os.Exit(1)
	}
}
