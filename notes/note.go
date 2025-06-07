package notes

import (
	"time"
)

type Note struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
