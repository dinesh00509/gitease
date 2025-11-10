package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinesh00509/GitEase/internals"
)

func main() {
	p := tea.NewProgram(internals.IntialModel())

	if err := p.Start(); err != nil {
		fmt.Println("Error running GitFlow:", err)
		os.Exit(1)
	}
}
