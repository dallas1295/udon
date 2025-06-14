package notes

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Store struct {
	notesDir string
}

// TODO Take a look at permissions and path for notes (I'm not a big fan of storage in Documents)
// Initialize the directory needed by udon to store notes.
func (s *Store) Init() error {
	// Get home directory
	usr, err := user.Current()
	if err != nil {
		return nil
	}

	// Set where Udon will store and retrieve notes locally
	// TODO change this for a configuration approach by the user
	s.notesDir = filepath.Join(usr.HomeDir, "Documents", "udon")

	// Create if does not exist
	if err := os.MkdirAll(s.notesDir, 0755); err != nil {
		fmt.Printf("There was an error creating note path: %s", err)
	}

	return nil
}

// Retrieve the notes from the storage directory
func (s *Store) GetNotes() error {
}

// Save note into the storage directory
func (s *Store) Save() error {
}

// Delete the note from the storage directory
func (s *Store) Delete() error {
}

func (s *Store) Update() error {
}
