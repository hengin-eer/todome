package todo

import (
	"testing"
	"time"
)

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func TestParseSimple(t *testing.T) {
	task := Parse("買い物に行く")
	if task.Done {
		t.Error("expected not done")
	}
	if task.Text != "買い物に行く" {
		t.Errorf("expected '買い物に行く', got %q", task.Text)
	}
}

func TestParseWithPriority(t *testing.T) {
	task := Parse("(A) 重要なタスク")
	if task.Priority != "A" {
		t.Errorf("expected priority A, got %q", task.Priority)
	}
	if task.Text != "重要なタスク" {
		t.Errorf("expected '重要なタスク', got %q", task.Text)
	}
}

func TestParseWithDate(t *testing.T) {
	task := Parse("2026-02-23 レポートを書く")
	if !task.CreatedAt.Equal(date(2026, 2, 23)) {
		t.Errorf("expected 2026-02-23, got %v", task.CreatedAt)
	}
	if task.Text != "レポートを書く" {
		t.Errorf("expected 'レポートを書く', got %q", task.Text)
	}
}

func TestParseDone(t *testing.T) {
	task := Parse("x 2026-02-23 2026-02-20 牛乳を買う +買い物")
	if !task.Done {
		t.Error("expected done")
	}
	if !task.CompletedAt.Equal(date(2026, 2, 23)) {
		t.Errorf("expected completed 2026-02-23, got %v", task.CompletedAt)
	}
	if !task.CreatedAt.Equal(date(2026, 2, 20)) {
		t.Errorf("expected created 2026-02-20, got %v", task.CreatedAt)
	}
	if len(task.Projects) != 1 || task.Projects[0] != "買い物" {
		t.Errorf("expected project '買い物', got %v", task.Projects)
	}
}

func TestParseProjectsAndContexts(t *testing.T) {
	task := Parse("2026-02-23 企画書を書く +仕事 @PC")
	if len(task.Projects) != 1 || task.Projects[0] != "仕事" {
		t.Errorf("expected project '仕事', got %v", task.Projects)
	}
	if len(task.Contexts) != 1 || task.Contexts[0] != "PC" {
		t.Errorf("expected context 'PC', got %v", task.Contexts)
	}
}

func TestSerializeSimple(t *testing.T) {
	task := Task{Text: "買い物に行く", CreatedAt: date(2026, 2, 23)}
	got := Serialize(task)
	expected := "2026-02-23 買い物に行く"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSerializeDone(t *testing.T) {
	task := Task{
		Done:        true,
		CompletedAt: date(2026, 2, 23),
		CreatedAt:   date(2026, 2, 20),
		Text:        "牛乳を買う +買い物",
	}
	got := Serialize(task)
	expected := "x 2026-02-23 2026-02-20 牛乳を買う +買い物"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSerializeWithPriority(t *testing.T) {
	task := Task{Priority: "A", CreatedAt: date(2026, 2, 23), Text: "重要なタスク"}
	got := Serialize(task)
	expected := "(A) 2026-02-23 重要なタスク"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRoundTrip(t *testing.T) {
	lines := []string{
		"(A) 2026-02-23 企画書を書く +仕事 @PC",
		"x 2026-02-23 2026-02-20 牛乳を買う +買い物",
		"2026-02-23 レポートを書く",
		"シンプルなタスク",
	}
	for _, line := range lines {
		task := Parse(line)
		got := Serialize(task)
		if got != line {
			t.Errorf("roundtrip failed:\n  input:  %q\n  output: %q", line, got)
		}
	}
}
