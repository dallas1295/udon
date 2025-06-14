package notes

import (
	"time"
)

type Note struct {
	Title   string
	Body    string
	ModTime time.Time
}
