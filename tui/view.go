package tui

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Gruvbox color palette
const (
	gruvBg     = "#282828"
	gruvFg     = "#ebdbb2"
	gruvYellow = "#fabd2f"
	gruvOrange = "#fe8019"
	gruvBlue   = "#83a598"
	gruvAqua   = "#8ec07c"
	gruvGreen  = "#b8bb26"
	gruvRed    = "#fb4934"
	gruvGray   = "#928374"
)

// listStyle defines the appearance of the notes list pane (left sidebar).
var listStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(gruvYellow)).
	Background(lipgloss.Color(gruvBg)).
	Foreground(lipgloss.Color(gruvFg)).
	Padding(1, 2)

// editorStyle defines the appearance of the editor/preview pane (main area).
var editorStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(gruvBlue)).
	Background(lipgloss.Color(gruvBg)).
	Foreground(lipgloss.Color(gruvFg)).
	Padding(1, 2)

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
		// fallback for startup if window size hasn't been set yet
		m.width = 80
	}
	listWidth := m.width / 3
	editorWidth := m.width - listWidth - 1

	// Build the notes list pane (left sidebar)
	var listContent strings.Builder
	for i, n := range m.notes {
		prefix := "  "
		if i == m.listIndex {
			prefix = "➤ "
		}
		listContent.WriteString(prefix + n.Title + "\n")
	}
	listPane := listStyle.Width(listWidth).Render(listContent.String())

	// Build the editor/preview pane (main area)
	var editorContent string
	switch m.state {
	case editorView:
		// Show the textarea for editing the note
		editorContent = m.textarea.View()
	default:
		// In listView (and any other state), show the content of the selected note rendered as Markdown
		if len(m.notes) > 0 {
			rendered, err := glamourRenderer.Render(m.notes[m.listIndex].Content)
			if err != nil {
				editorContent = m.notes[m.listIndex].Content // fallback to plain text
			} else {
				editorContent = rendered
			}
		} else {
			editorContent = "Select a note to preview/edit."
		}
	}
	editorPane := editorStyle.Width(editorWidth).Render(editorContent)

	// Join the two panes horizontally to create the split layout
	return lipgloss.JoinHorizontal(lipgloss.Top, listPane, editorPane)
}
