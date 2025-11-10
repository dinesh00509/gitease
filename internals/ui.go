package internals

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type step struct {
	label string
	done  bool
}

type model struct {
	cursor  int
	steps   []step
	output  string
	running bool
}

func IntialModel() model {
	return model{
		cursor: 0,
		steps: []step{
			{"stage all the changes", false},
			{"show git status", false},
			{"commit the changes", false},
			{"Push to remote", false},
		},
	}
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString("GIT FLOW\n")
	b.WriteString("-------------------------\n")

	for i, s := range m.steps {
		var cursor = " "
		if m.cursor == i {
			cursor = "->"
		}
		status := "[ ]"
		if s.done {
			status = "[✔]"
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, status, s.label))
	}
	b.WriteString("\nUse ↑/↓ to navigate, Enter to run, q to quit.\n")
	b.WriteString("───────────────────────────────\n")
	b.WriteString(m.output)

	return b.String()
}

func (m model) runCurrentStep() (tea.Model, tea.Cmd) {
	step := m.steps[m.cursor]
	m.output = fmt.Sprintf("Running: %s...\n", m.steps[m.cursor].label)

	switch m.cursor {
	case 0:
		m.output += RunGit("add", ".")
	case 1:
		m.output += RunGit("status")
	case 2:
		m.output += RunGit("commit", "-m", "Auto commit from GitFlow")
	case 3:
		m.output += RunGit("push")
	}
	step.done = true
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.steps)-1 {
				m.cursor++
			}

		case "enter":
			return m.runCurrentStep()
		}
	}
	return m, nil
}
