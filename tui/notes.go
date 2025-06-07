package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FIXME this is borrowed code from a previous project, need to make sure the bubbletea structure will be the same before commiting to this model approach
type model struct {
	state uint
	// store     *store
	// notes     []Note
	// currNote  Note
	listIndex int
	textarea  textarea.Model
	textinput textinput.Model
}

func NewModel(store *Store) model {
}

func (m model) Init() tea.Cmd {
	return nil
}
