package internals

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const asciiArt = `         _               _        _            _            _                  _            _       
        /\ \            /\ \     /\ \         /\ \         / /\               / /\         /\ \     
       /  \ \           \ \ \    \_\ \       /  \ \       / /  \             / /  \       /  \ \    
      / /\ \_\          /\ \_\   /\__ \     / /\ \ \     / / /\ \           / / /\ \__   / /\ \ \   
     / / /\/_/         / /\/_/  / /_ \ \   / / /\ \_\   / / /\ \ \         / / /\ \___\ / / /\ \_\  
    / / / ______      / / /    / / /\ \ \ / /_/_ \/_/  / / /  \ \ \        \ \ \ \/___// /_/_ \/_/  
   / / / /\_____\    / / /    / / /  \/_// /____/\    / / /___/ /\ \        \ \ \     / /____/\     
  / / /  \/____ /   / / /    / / /      / /\____\/   / / /_____/ /\ \   _    \ \ \   / /\____\/     
 / / /_____/ / /___/ / /__  / / /      / / /______  / /_________/\ \ \ /_/\__/ / /  / / /______     
/ / /______\/ //\__\/_/___\/_/ /      / / /_______\/ / /_       __\ \_\\ \/___/ /  / / /_______\    
\/___________/ \/_________/\_\/       \/__________/\_\___\     /____/_/ \_____\/   \/__________/    `

var (
	asciiArtStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)

	// titleStyle = lipgloss.NewStyle().
	// 		Bold(true).
	// 		Foreground(lipgloss.Color("#00FFFF")).
	// 		Background(lipgloss.Color("#1E1E1E")).
	// 		Padding(0, 2).
	// 		MarginBottom(1).
	// 		Underline(true)

	subtleDivider = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render(strings.Repeat("─", 50))

	stepPendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	stepActiveStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	stepDoneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Bold(true)

	outputBox = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#555555")).Padding(1, 2).MarginTop(1).Width(65)

	successText = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Bold(true)
	errorText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	hintText    = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Faint(true).MarginTop(1)
)

func (m Model) View() string {
	var content strings.Builder

	content.WriteString(asciiArtStyle.Render(asciiArt) + "\n\n")
	content.WriteString(subtleDivider + "\n\n")

	for i, s := range m.Steps {
		var line string
		if s.Done {
			line = stepDoneStyle.Render(fmt.Sprintf("✔ %s", s.Label))
		} else if i == m.Cursor {
			line = stepActiveStyle.Render(fmt.Sprintf("▶ %s", s.Label))
		} else {
			line = stepPendingStyle.Render(fmt.Sprintf("• %s", s.Label))
		}
		content.WriteString(line + "\n")
	}
	content.WriteString("\n" + subtleDivider + "\n\n")

	if m.Committing {
		content.WriteString(stepActiveStyle.Render("Commit Message:") + "\n")
		content.WriteString(m.TextInput.View())
	} else if m.BranchMode {
		if m.NewBranch {
			content.WriteString(stepActiveStyle.Render("Enter new branch name:") + "\n")
		} else {
			content.WriteString(stepActiveStyle.Render("Enter branch name to switch:") + "\n")
		}
		content.WriteString(m.TextInput.View())
	} else if m.PullBranch {
		if m.PullFromOtherBranch {
			content.WriteString(stepActiveStyle.Render("Enter branch name to pull from:") + "\n")
			content.WriteString(m.TextInput.View())
		} else {
			content.WriteString(stepActiveStyle.Render("Press Enter to pull from current branch") + "\n")
			content.WriteString(m.TextInput.View())
		}
	} else {
		if m.Output == "" {
			content.WriteString(outputBox.Render("Waiting for command..."))
		} else if strings.Contains(strings.ToLower(m.Output), "error") {
			content.WriteString(outputBox.Render(errorText.Render(m.Output)))
		} else if strings.Contains(strings.ToLower(m.Output), "completed") ||
			strings.Contains(strings.ToLower(m.Output), "success") {
			content.WriteString(outputBox.Render(successText.Render(m.Output)))
		} else {
			content.WriteString(outputBox.Render(m.Output))
		}
	}

	content.WriteString("\n\n" + subtleDivider + "\n")
	content.WriteString(hintText.Render("↑/↓ navigate • Enter run • q quit • ESC cancel input") + "\n")

	contentStr := content.String()
	centered := lipgloss.Place(
		lipgloss.Width(contentStr),
		0,
		lipgloss.Center,
		lipgloss.Center,
		contentStr,
	)
	return centered

}
