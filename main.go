package main

import (
	"fmt"
	"gitflow/internals"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(internals.IntialModel())

	if err := p.Start(); err != nil {
		fmt.Println("Error running GitFlow:", err)
		os.Exit(1)
	}
}
