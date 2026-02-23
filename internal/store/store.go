package store

import "github.com/hengin-eer/todome/internal/todo"

// Store defines the interface for task persistence.
type Store interface {
	Load() ([]todo.Task, error)
	Save(tasks []todo.Task) error
}
