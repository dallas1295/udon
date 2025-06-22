// Package notes provides storage and business logic for plain text notes.
package notes

import "strings"

// LoadToMem loads the note with the given title from the Store and returns it for use in memory.
// It returns an error if the note cannot be loaded.
func LoadToMem(store *Store, title string) (*Note, error) {
	note, err := store.Load(title)
	if err != nil {
		return nil, err
	}
	return note, nil
}

// SearchNotes returns a slice of notes whose title or content contains the
// given query string, case-insensitively. It retrieves all notes from the
// provided Store and filters them based on the query. If an error occurs
// while retrieving notes, it returns nil and the error.
func SearchNotes(store *Store, query string) ([]Note, error) {
	notes, err := store.GetNotes()
	if err != nil {
		return nil, err
	}

	var results []Note
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Title), strings.ToLower(query)) || strings.Contains(strings.ToLower(note.Content), strings.ToLower(query)) {
			results = append(results, note)
		}
	}
	return results, nil
}

// checkDuplicateName checks if a note with the given title already exists in the Store.
// It returns "duplicate" if a note exists, "okay" otherwise, or an error if retrieval fails.
func checkDuplicateName(store *Store, title string) (string, error) {
	notes, err := store.GetNotes()
	if err != nil {
		return "", err
	}

	for _, note := range notes {
		if note.Title == title {
			return "duplicate", nil
		}
	}

	return "okay", nil
}

// ConfirmSave saves the given note using the provided Store if confirm is true and
// there are no preexisting notes with the same name. It returns a status string
// indicating the result ("saved", "note exists", or "not saved") and an error if
// saving fails or a duplicate is found.
func ConfirmSave(store *Store, note Note, confirm bool) (string, error) {
	if confirm {
		status, err := checkDuplicateName(store, note.Title)
		if err != nil {
			return "", err
		}
		if status == "duplicate" {
			return "note exists", nil
		}
		if err := store.Save(note); err != nil {
			return "", err
		}
		return "saved", nil
	}
	return "not saved", nil
}

// ConfirmClose saves the given note using the provided Store if confirm is true.
// It returns "saved" if the note was saved, "unsaved" otherwise, and any error encountered.
func ConfirmClose(store *Store, note Note, confirm bool) (string, error) {
	if confirm {
		if err := store.Save(note); err != nil {
			return "", err
		}
		return "saved", nil
	}
	return "unsaved", nil
}

// ConfirmUpdate updates the note's title and/or content using the provided Store
// if confirm is true and there are no duplicate titles (if the title is changing).
// It returns a status string indicating the result ("updated", "duplicate", "not updated")
// and any error encountered.
func ConfirmUpdate(store *Store, oldTitle string, updatedTitle *string, updatedContent *string, confirm bool) (string, error) {
	if !confirm {
		return "not updated", nil
	}

	// Check if the title is being changed and if so, check for duplicates
	if updatedTitle != nil && strings.TrimSpace(*updatedTitle) != "" && *updatedTitle != oldTitle {
		status, err := checkDuplicateName(store, *updatedTitle)
		if err != nil {
			return "", err
		}
		if status == "duplicate" {
			return "duplicate", nil
		}
	}

	// Call the storage layer's Update function
	if err := store.Update(oldTitle, updatedTitle, updatedContent); err != nil {
		return "", err
	}

	return "updated", nil
}

// ConfirmDelete deletes the note with the given title from the Store if confirm is true.
// It returns "deleted" if the note was deleted, "not deleted" otherwise, and any error encountered.
func ConfirmDelete(store *Store, note string, confirm bool) (string, error) {
	if confirm {
		if err := store.Delete(note); err != nil {
			return "", err
		}
		return "deleted", nil
	}
	return "not deleted", nil
}
