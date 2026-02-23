package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hengin-eer/todome/internal/todo"
)

func TestFileStoreLoadEmpty(t *testing.T) {
	s := NewFileStore("/tmp/nonexistent-todome-test.txt")
	tasks, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestFileStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.txt")
	s := NewFileStore(path)

	tasks := []todo.Task{
		{CreatedAt: time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC), Text: "タスク1 +仕事"},
		{Done: true, CompletedAt: time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC), CreatedAt: time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC), Text: "タスク2"},
	}

	if err := s.Save(tasks); err != nil {
		t.Fatal(err)
	}

	// Verify file content
	data, _ := os.ReadFile(path)
	content := string(data)
	if content != "2026-02-23 タスク1 +仕事\nx 2026-02-23 2026-02-20 タスク2\n" {
		t.Errorf("unexpected file content:\n%s", content)
	}

	// Load back
	loaded, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Text != "タスク1 +仕事" {
		t.Errorf("expected 'タスク1 +仕事', got %q", loaded[0].Text)
	}
	if !loaded[1].Done {
		t.Error("expected task 2 to be done")
	}
}
