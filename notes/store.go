package notes

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
func (s *Store) GetNotes() ([]Note, error) {
	entries, err := os.ReadDir(s.notesDir)
	if err != nil {
		fmt.Printf("Error retrieving notes from directory: %v", err)
		return nil, err
	}

	var notes []Note

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		var note Note
		var content []string

		// Get note modified time from file for sorting
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", entry.Name(), err) // output the erroneous file rea
			continue                                                     // skipping unreadable files
		}

		modTime := info.ModTime().Local()

		// Get note name from filename
		filename := entry.Name()
		name := strings.TrimSuffix(filename, ".txt")

		// Get note content
		path := filepath.Join(s.notesDir, entry.Name())
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filename, err) // output which files cannot be opened
			continue                                                 // skipping unreadable files
		}

		// Get note content from inside the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			content = append(content, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err) // output the erroneous file read
			file.Close()
			continue
		}

		// Output information for our function
		note.ModTime = modTime
		note.Title = name
		note.Body = strings.Join(content, "\n")
		notes = append(notes, note)

		// Close the files
		file.Close()
	}

	return notes, nil
}

// Save note into the storage directory
func (s *Store) Save() error {
}

// Delete the note from the storage directory
func (s *Store) Delete() error {
}

func (s *Store) Update() error {
}
