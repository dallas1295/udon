package notes

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type Note struct {
	Title   string
	Content string
	ModTime time.Time
}

type Store struct {
	notesDir string
}

// Init initializes the storage directory. If the directory does not exist, it creates one.
func (s *Store) Init() error {
	// Get the current user's home directory.
	usr, err := user.Current()
	if err != nil {
		return nil
	}

	// Set the path where notes will be stored.
	// TODO: Allow user configuration for storage location.
	s.notesDir = filepath.Join(usr.HomeDir, "Documents", "udon")

	// Create the notes directory if it does not exist.
	if err := os.MkdirAll(s.notesDir, 0755); err != nil {
		fmt.Printf("There was an error creating note path: %s", err)
	}

	return nil
}

// GetNotes retrieves all notes from the storage directory.
// It returns a slice of Note and any error encountered during retrieval.
func (s *Store) GetNotes() ([]Note, error) {
	// Read all entries in the notes directory.
	entries, err := os.ReadDir(s.notesDir)
	if err != nil {
		fmt.Printf("Error retrieving notes from directory: %v", err)
		return nil, err
	}

	var notes []Note

	for _, entry := range entries {
		// Skip directories; only process files.
		if entry.IsDir() {
			continue
		}

		var note Note
		var content []string

		// Retrieve file metadata for modification time.
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Error reading file info for %s: %v\n", entry.Name(), err)
			continue // Skip files with unreadable metadata.
		}
		modTime := info.ModTime().Local()

		// Extract the note title from the filename.
		filename := entry.Name()
		name := strings.TrimSuffix(filename, ".txt")

		// Build the full path to the note file.
		path := filepath.Join(s.notesDir, entry.Name())

		// Open the note file for reading.
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filename, err)
			continue // Skip files that cannot be opened.
		}

		// Read the note content line by line.
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			content = append(content, scanner.Text())
		}
		// Handle any scanning errors.
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			file.Close()
			continue
		}

		// Populate the Note struct with file data.
		note.ModTime = modTime
		note.Title = name
		note.Content = strings.Join(content, "\n")
		notes = append(notes, note)

		// Close the file after reading.
		file.Close()
	}

	return notes, nil
}

// Save writes the given note to the storage directory.
// It trims trailing whitespace from the note content before saving.
func (s *Store) Save(note Note) error {
	// Build the filename from the note's title.
	filename := strings.TrimSpace(note.Title) + ".txt"

	// Remove trailing whitespace from the note content.
	content := strings.TrimRightFunc(note.Content, unicode.IsSpace)

	// Build the full path for the note file.
	path := filepath.Join(s.notesDir, filename)

	// Write the note content to the file.
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Could not save note w/ error: %v", err)
	}

	return nil
}

// Delete removes the note with the specified name from the storage directory.
// If the note does not exist, Delete prints a message and returns an error.
func (s *Store) Delete(noteName string) error {
	// Build the full path to the note file.
	filename := strings.TrimSpace(noteName) + ".txt"
	toDelete := filepath.Join(s.notesDir, filename)

	// Attempt to remove the note file.
	err := os.Remove(toDelete)
	if err != nil {
		// If the file does not exist, print a message.
		if os.IsNotExist(err) {
			fmt.Printf("A note with the name %s doesn't exist\n", noteName)
		}
		// Print any other error encountered during deletion.
		fmt.Printf("There was an error deleting this note: %v\n", err)
		return err
	}

	return nil
}

// Update modifies the title and/or content of an existing note.
// If both updatedTitle and updatedContent are nil, Update returns without making changes.
func (s *Store) Update(oldTitle string, updatedTitle *string, updatedContent *string) error {
	// Return early if there are no updates to apply.
	if updatedTitle == nil && updatedContent == nil {
		fmt.Printf("There are no updates for %s\n", oldTitle)
		return nil
	}

	// Build the path to the current note file.
	filename := strings.TrimSpace(oldTitle) + ".txt"
	path := filepath.Join(s.notesDir, filename)

	// Rename the note file if a new, non-empty title is provided.
	if updatedTitle != nil && strings.TrimSpace(*updatedTitle) != "" {
		updatedFilename := strings.TrimSpace(*updatedTitle) + ".txt"
		newPath := filepath.Join(s.notesDir, updatedFilename)
		err := os.Rename(path, newPath)
		if err != nil {
			fmt.Printf("There was an error renaming %s to %s: %v\n", oldTitle, *updatedTitle, err)
			return err
		}
		// Update the path to point to the renamed file.
		path = newPath
	}

	// Update the note content if new content is provided.
	if updatedContent != nil {
		oldBytes, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("There was an error reading %s content: %v\n", path, err)
			return err
		}
		oldContent := string(oldBytes)
		newContent := strings.TrimRightFunc(*updatedContent, unicode.IsSpace)

		// Only write if the new content is different from the old content.
		if newContent != oldContent {
			err := os.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("There was an error writing new content to %s: %v\n", path, err)
				return err
			}
		}
	}

	return nil
}
