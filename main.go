/*
Copyright © 2025 Dallas S.
*/
package main

import (
	"log"
	"udon/notes"
	"udon/tui"
)

func main() {
	store := &notes.Store{}
	if err := store.Init(); err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	if err := tui.Run(store); err != nil {
		log.Fatalf("TUI exited with error: %v", err)
	}
}
