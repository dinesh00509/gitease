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
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFFF")).
			Background(lipgloss.Color("#1E1E1E")).
			Padding(0, 2).
			MarginBottom(1).
			Underline(true)

	subtleDivider = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			Render(strings.Repeat("─", 50))

	stepPendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	stepActiveStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	stepDoneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Bold(true)

	outputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Padding(1, 2).
			MarginTop(1).
			Width(65)

	successText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF7F")).
			Bold(true)

	errorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	hintText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Faint(true).
			MarginTop(1)
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
	ti.Placeholder = "Enter your commit message..."
	ti.CharLimit = 120
	ti.Width = 50

	return model{
		cursor: 0,
		steps: []step{
			{"Stage all changes", false},
			{"Show Git status", false},
			{"Commit changes", false},
			{"Push to remote", false},
		},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🚀 GitEase – Smart Git CLI Assistant") + "\n")
	b.WriteString(subtleDivider + "\n\n")

	for i, s := range m.steps {
		var line string

		if s.done {
			line = stepDoneStyle.Render(fmt.Sprintf("✔ %s", s.label))
		} else if i == m.cursor {
			line = stepActiveStyle.Render(fmt.Sprintf("▶ %s", s.label))
		} else {
			line = stepPendingStyle.Render(fmt.Sprintf("• %s", s.label))
		}

		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + subtleDivider + "\n\n")

	if m.commiting {
		b.WriteString(stepActiveStyle.Render("💬 Commit Message:") + "\n")
		b.WriteString(m.textInput.View())
	} else {
		if m.output == "" {
			b.WriteString(outputBox.Render("Waiting for command..."))
		} else if strings.Contains(strings.ToLower(m.output), "error") {
			b.WriteString(outputBox.Render(errorText.Render(m.output)))
		} else if strings.Contains(strings.ToLower(m.output), "completed") ||
			strings.Contains(strings.ToLower(m.output), "success") {
			b.WriteString(outputBox.Render(successText.Render(m.output)))
		} else {
			b.WriteString(outputBox.Render(m.output))
		}
	}

	b.WriteString("\n\n" + subtleDivider + "\n")
	b.WriteString(hintText.Render("↑/↓ navigate • Enter run • q quit • ESC cancel commit") + "\n")

	return b.String()
}

func (m model) runCurrentStep() (tea.Model, tea.Cmd) {
	m.output = fmt.Sprintf("Running: %s...\n", m.steps[m.cursor].label)
	switch m.cursor {
	case 0:
		m.output += RunGit("add", ".")
		m.output += "\n✅ Staging completed successfully."

	case 1:
		m.output += RunGit("status")
		m.output += "\n✅ Status check completed."

	case 2:
		m.commiting = true
		m.textInput.SetValue("")
		m.textInput.Focus()
		m.output = ""

	case 3:
		m.output += RunGit("push")
		m.output += "\n🚀 Push completed successfully.\nAll steps done — exiting..."
		m.steps[m.cursor].done = true
		return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg { return tea.QuitMsg{} })
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
					m.output = errorText.Render("❌ Commit message cannot be empty.")
				}
				m.commiting = false
				return m, nil
			case "esc":
				m.commiting = false
				m.output = "Commit action cancelled."
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
