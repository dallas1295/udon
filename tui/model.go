package tui

import (
	"log"
	"udon/notes"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	editorView
	previewView
)

type model struct {
	state     uint
	store     *notes.Store
	notes     []notes.Note
	currNote  notes.Note
	statusMsg string
	listIndex int
	textarea  textarea.Model
	textinput textinput.Model
	width     int
	height    int
}

func (m model) Init() tea.Cmd {
	return nil
}

func NewModel(store *notes.Store) model {
	notesSlice, err := store.GetNotes()
	if err != nil {
		log.Fatalf("Error retrieving notes: %v", err)
	}
	return model{
		state:     listView,
		store:     store,
		notes:     notesSlice,
		textarea:  textarea.New(),
		textinput: textinput.New(),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case editorView:
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = listView
			}
		case listView:
			switch key {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "k", "up":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "j", "down":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter", "l":
				notePtr, err := notes.LoadToMem(m.store, m.notes[m.listIndex].Title)
				if err != nil {
					m.statusMsg = "Error loading note: " + err.Error()
				} else if notePtr != nil {
					m.currNote = *notePtr
					m.textarea.SetValue(m.currNote.Content)
					m.state = editorView
					m.textarea.Focus()
				}
			}
		case previewView:
			switch key {
			case "e":
				m.textarea.SetValue(m.currNote.Content)
				m.state = editorView
				m.textarea.Focus()
			}
		}
	}
	return m, tea.Batch(cmds...)
}
