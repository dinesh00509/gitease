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
			fmt.Println("Usage: gitease")
			fmt.Println("Navigate and manage Git tasks interactively in your terminal.")
			fmt.Println("Flags:")
			fmt.Println("  --version, -v   Show version")
			fmt.Println("  --help, -h      Show this help message")
			os.Exit(0)
		}
		if arg == "--run" || arg == "-r" {
			if len(os.Args) < 3 {
				fmt.Println("Error: Missing argument for --run")
				fmt.Println("Usage: gitease --run <command>")
				os.Exit(1)
			}
			command := os.Args[2]
			output := internals.RunGit(command)
			fmt.Println(output)
			os.Exit(0)
		}
	}

	p := tea.NewProgram(internals.InitialModel())

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running GitEase:", err)
		os.Exit(1)
	}
}
