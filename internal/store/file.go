package store

import (
	"bufio"
	"os"
	"strings"

	"github.com/hengin-eer/todome/internal/todo"
)

// FileStore reads and writes tasks to a todo.txt file.
type FileStore struct {
	Path string
}

// NewFileStore creates a new FileStore for the given path.
func NewFileStore(path string) *FileStore {
	return &FileStore{Path: path}
}

// Load reads all tasks from the file. Returns empty slice if file doesn't exist.
func (fs *FileStore) Load() ([]todo.Task, error) {
	f, err := os.Open(fs.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []todo.Task{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var tasks []todo.Task
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tasks = append(tasks, todo.Parse(line))
	}
	return tasks, scanner.Err()
}

// Save writes all tasks to the file, overwriting existing content.
func (fs *FileStore) Save(tasks []todo.Task) error {
	f, err := os.Create(fs.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, t := range tasks {
		w.WriteString(todo.Serialize(t))
		w.WriteString("\n")
	}
	return w.Flush()
}
