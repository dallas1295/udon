package notes

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

// Note represents a single note with a title, content, and modification time.
type Note struct {
	Title   string
	Content string
	ModTime time.Time
}

// Store manages the storage and retrieval of notes from the filesystem.
type Store struct {
	notesDir string
}

// sanitizeFilename replaces invalid filename characters with underscores.
func sanitizeFilename(name string) string {
	// invalidFilenameChars is a regexp for characters not allowed in filenames.
	invalidFilenameChars := regexp.MustCompile(`[<>:"/\\|?*]`)

	return invalidFilenameChars.ReplaceAllString(name, "_")
}

// Init initializes the storage directory. If the directory does not exist, it creates one.
func (s *Store) Init() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}

	s.notesDir = filepath.Join(usr.HomeDir, "Documents", "udon")

	if err := os.MkdirAll(s.notesDir, 0755); err != nil {
		return fmt.Errorf("error creating note path: %w", err)
	}

	return nil
}

// GetNotes retrieves all notes from the storage directory.
// It returns a slice of Note and any error encountered during retrieval.
func (s *Store) GetNotes() ([]Note, error) {
	entries, err := os.ReadDir(s.notesDir)
	if err != nil {
		return nil, fmt.Errorf("error reading notes directory: %w", err)
	}

	var notes []Note

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) != ".md" {
			continue // Only process .md files
		}

		name := strings.TrimSuffix(filename, ".md")
		path := filepath.Join(s.notesDir, filename)

		file, err := os.Open(path)
		if err != nil {
			// Skip files that cannot be opened, but continue processing others
			continue
		}

		var content []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			content = append(content, scanner.Text())
		}
		file.Close()

		info, err := entry.Info()
		if err != nil {
			continue // Skip files with unreadable metadata
		}

		note := Note{
			Title:   name,
			Content: strings.Join(content, "\n"),
			ModTime: info.ModTime().Local(),
		}
		notes = append(notes, note)
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ModTime.After(notes[j].ModTime)
	})

	return notes, nil
}

// Load retrieves a single note by title.
// Returns a pointer to the Note and any error encountered.
func (s *Store) Load(title string) (*Note, error) {
	filename := sanitizeFilename(strings.TrimSpace(title)) + ".md"
	path := filepath.Join(s.notesDir, filename)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("note %q does not exist", title)
		}
		return nil, fmt.Errorf("error reading note %q: %w", title, err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error getting file info for %q: %w", title, err)
	}

	return &Note{
		Title:   title,
		Content: string(content),
		ModTime: info.ModTime().Local(),
	}, nil
}

// Save writes the given note to the storage directory.
// It trims trailing whitespace from the note content before saving.
func (s *Store) Save(note Note) error {
	if strings.TrimSpace(note.Title) == "" {
		return errors.New("note title cannot be empty")
	}
	filename := sanitizeFilename(strings.TrimSpace(note.Title)) + ".md"
	content := strings.TrimRightFunc(note.Content, unicode.IsSpace)
	path := filepath.Join(s.notesDir, filename)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("could not save note: %w", err)
	}
	return nil
}

// Delete removes the note with the specified name from the storage directory.
// If the note does not exist, Delete returns an error.
func (s *Store) Delete(noteName string) error {
	if strings.TrimSpace(noteName) == "" {
		return errors.New("note name cannot be empty")
	}
	filename := sanitizeFilename(strings.TrimSpace(noteName)) + ".md"
	toDelete := filepath.Join(s.notesDir, filename)

	err := os.Remove(toDelete)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("note %q does not exist", noteName)
		}
		return fmt.Errorf("error deleting note %q: %w", noteName, err)
	}
	return nil
}

// Update modifies the title and/or content of an existing note.
// If both updatedTitle and updatedContent are nil, Update returns without making changes.
func (s *Store) Update(oldTitle string, updatedTitle *string, updatedContent *string) error {
	if strings.TrimSpace(oldTitle) == "" {
		return errors.New("old title cannot be empty")
	}
	if updatedTitle == nil && updatedContent == nil {
		return nil // Nothing to update
	}

	oldFilename := sanitizeFilename(strings.TrimSpace(oldTitle)) + ".md"
	oldPath := filepath.Join(s.notesDir, oldFilename)
	path := oldPath

	// Rename the note file if a new, non-empty title is provided.
	if updatedTitle != nil && strings.TrimSpace(*updatedTitle) != "" {
		newFilename := sanitizeFilename(strings.TrimSpace(*updatedTitle)) + ".md"
		newPath := filepath.Join(s.notesDir, newFilename)
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("error renaming %q to %q: %w", oldTitle, *updatedTitle, err)
		}
		path = newPath
	}

	// Update the note content if new content is provided.
	if updatedContent != nil {
		oldBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading %q content: %w", path, err)
		}
		oldContent := string(oldBytes)
		newContent := strings.TrimRightFunc(*updatedContent, unicode.IsSpace)

		if newContent != oldContent {
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("error writing new content to %q: %w", path, err)
			}
		}
	}

	return nil
}
