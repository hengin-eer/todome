package todo

import "time"

// Task represents a single todo.txt task.
type Task struct {
	Done        bool
	Priority    string // "" or "A"-"Z"
	CompletedAt time.Time
	CreatedAt   time.Time
	Text        string
	Projects    []string // +project tags
	Contexts    []string // @context tags
}
