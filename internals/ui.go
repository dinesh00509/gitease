package internals

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	//colors
	green     = lipgloss.Color("#00FF00")
	yellow    = lipgloss.Color("#FFD700")
	red       = lipgloss.Color("#FF5555")
	cyan      = lipgloss.Color("#00FFFF")
	darkGray  = lipgloss.Color("#2E2E2E")
	lightGray = lipgloss.Color("#A9A9A9")

	titleBox = lipgloss.NewStyle().Bold(true).Foreground(cyan).Border(lipgloss.RoundedBorder()).BorderForeground(cyan).Padding(0, 2).Align(lipgloss.Center)

	doneStyle   = lipgloss.NewStyle().Foreground(green)
	cursorStyle = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	labelStyle  = lipgloss.NewStyle().Foreground(lightGray)

	outputBox = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lightGray).Padding(1, 2).Foreground(green)

	successStyle = lipgloss.NewStyle().Foreground(green).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(red).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(yellow)

	footerStyle = lipgloss.NewStyle().Foreground(lightGray).MarginTop(1).Faint(true)
)

type step struct {
	label string
	done  bool
}

type model struct {
	cursor    int
	steps     []step
	output    string
	commiting bool
	textInput textinput.Model
}

func IntialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your Message to Commit...."
	ti.Focus()
	ti.CharLimit = 120
	ti.Width = 40
	return model{
		cursor: 0,
		steps: []step{
			{"stage all the changes", false},
			{"show git status", false},
			{"commit the changes", false},
			{"Push to remote", false},
		},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

// For viewing.
func (m model) View() string {
	var b strings.Builder

	title := titleBox.Render("🧠 GIT FLOW ASSISTANT")
	b.WriteString(title + "\n\n")

	for i, s := range m.steps {
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("➤")
		}

		status := "[ ]"
		if s.done {
			status = doneStyle.Render("[✔]")
		} else {
			status = labelStyle.Render(status)
		}

		label := labelStyle.Render(s.label)
		line := fmt.Sprintf("%s %s %s\n", cursor, status, label)
		b.WriteString(line)
	}

	b.WriteString("\n" + lipgloss.NewStyle().
		Foreground(lightGray).
		Render(strings.Repeat("─", 40)) + "\n")

	if m.commiting {
		msg := infoStyle.Render("✏️  Type your commit message and press Enter (ESC to cancel):")
		b.WriteString(msg + "\n")
		b.WriteString(outputBox.Render(m.textInput.View()))
	} else {
		if strings.Contains(m.output, "Error") {
			b.WriteString(outputBox.Foreground(red).Render(m.output))
		} else if strings.Contains(m.output, "completed") {
			b.WriteString(outputBox.Foreground(green).Render(m.output))
		} else {
			b.WriteString(outputBox.Foreground(lightGray).Render(m.output))
		}
	}

	b.WriteString("\n" + lipgloss.NewStyle().
		Render(strings.Repeat("─", 40)) + "\n")
	b.WriteString(footerStyle.Render("↑↓ Navigate   ⏎ Run   q Quit"))

	return b.String()
}

func (m model) runCurrentStep() (tea.Model, tea.Cmd) {
	m.output = fmt.Sprintf("Running: %s...\n", m.steps[m.cursor].label)
	switch m.cursor {
	case 0:
		m.output += RunGit("add", ".")
		m.output += "\nStaging completed successfully."

	case 1:
		m.output += RunGit("status")
		m.output += "Status check completed."

	case 2:
		m.commiting = true
		m.textInput.SetValue("")
		m.textInput.Focus()
		m.output = ""

	case 3:
		m.output += RunGit("push")
		m.output += "Pushed to remote branch successfully."
		m.output += "\nAll steps completed. Exiting GitFlow..."
		m.steps[m.cursor].done = true

		return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
			return tea.QuitMsg{}
		})
	}

	if m.cursor != 2 {
		m.steps[m.cursor].done = true
	}

	return m, tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg { return nil })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.commiting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				commitMsg := strings.TrimSpace(m.textInput.Value())
				if commitMsg != "" {
					m.output = RunGit("commit", "-m", commitMsg)
					m.steps[m.cursor].done = true
				} else {
					m.output = "Commit message cannot be empty."
				}
				m.commiting = false
				return m, nil
			case "esc":
				m.commiting = false
				m.output = "Commit Action Cancelled"
				return m, nil
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

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
