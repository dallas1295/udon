package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// statusStyle defines the appearance of the status pane (top bar)
var statusStyle = lipgloss.NewStyle().
	Background(lipgloss.Color(gruvYellow)).
	Foreground(lipgloss.Color(gruvBg)).
	MarginTop(1).
	MarginRight(2).
	MarginLeft(2)

// listStyle defines the appearance of the notes list pane (left sidebar).
var listStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(gruvGreen)).
	Foreground(lipgloss.Color(gruvFg)).
	Padding(2, 1).
	MarginLeft(1)

// editorStyle defines the appearance of the editor/preview pane (main area).
var editorStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(gruvGreen)).
	Foreground(lipgloss.Color(gruvFg)).
	Padding(2, 1).
	MarginRight(1)

// glamourRenderer is a reusable Glamour renderer with Gruvbox theming.
var glamourRenderer = func() *glamour.TermRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(80),
	)
	return r
}()

// View renders the main TUI layout with a split-pane design.
// The left third of the terminal is the notes list, and the right two-thirds
// is the editor or preview pane. The layout is responsive to terminal resizing.
func (m model) View() string {
	if m.width == 0 {
		m.width = 80
	}
	listWidth := m.width/4 - 6
	editorWidth := m.width - listWidth - editorStyle.GetHorizontalFrameSize() - 1
	statusWidth := m.width - statusStyle.GetHorizontalFrameSize()

	var listContent strings.Builder
	for i, n := range m.notes {
		prefix := "  "
		if i == m.listIndex {
			prefix = "➤ "
		}
		listContent.WriteString(prefix + n.Title + "\n")
	}
	listPane := listStyle.Width(listWidth).Render(listContent.String())

	var editorContent string
	switch m.state {
	case editorView:
		editorContent = m.textarea.View()
	default:
		if len(m.notes) > 0 {
			rendered, err := glamourRenderer.Render(m.notes[m.listIndex].Content)
			if err != nil {
				editorContent = m.notes[m.listIndex].Content
			} else {
				editorContent = rendered
			}
		} else {
			editorContent = "Select a note to preview/edit."
		}
	}
	editorPane := editorStyle.Width(editorWidth).Render(editorContent)

	statusLeft := "Status"
	var rightStatus string
	if m.state == editorView {
		rightStatus = fmt.Sprintf("%s - Editing", m.currNote.Title)
	} else {
		rightStatus = fmt.Sprintf("%s - Preview", m.currNote.Title)
	}
	statusContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusLeft,
		lipgloss.PlaceHorizontal(
			statusWidth-lipgloss.Width(statusLeft),
			lipgloss.Right,
			rightStatus,
		),
	)
	statusPane := statusStyle.Width(statusWidth).Render(statusContent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusPane,
		lipgloss.JoinHorizontal(lipgloss.Top, listPane, editorPane),
	)
}
