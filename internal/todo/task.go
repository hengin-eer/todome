package todo

import "time"

// Task represents a single todo.txt task.
type Task struct {
	Done             bool
	Priority         string // "" or "A"-"Z"
	CompletedAt      time.Time
	CreatedAt        time.Time
	DueDate          time.Time // due:YYYY-MM-DD or due:YYYY-MM-DDTHH:mm:ss
	Text             string
	Projects         []string // +project tags
	Contexts         []string // @context tags
	Note             string   // completion note (note:...)
	CreatedHasTime   bool     // true if CreatedAt includes HH:mm:ss
	CompletedHasTime bool     // true if CompletedAt includes HH:mm:ss
	DueHasTime       bool     // true if DueDate includes HH:mm:ss
}
