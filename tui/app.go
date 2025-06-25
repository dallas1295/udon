// Package tui provides a terminal user interface for managing and editing notes.
package tui

import (
	"udon/notes"

	tea "github.com/charmbracelet/bubbletea"
)

// Run sets up and starts the Bubbletea TUI program.
// It returns any error encountered during execution.
func Run(store *notes.Store) error {
	m := NewModel(store)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
