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
		"2026-02-23 請求書処理 +仕事 due:2026-03-01",
	}
	for _, line := range lines {
		task := Parse(line)
		got := Serialize(task)
		if got != line {
			t.Errorf("roundtrip failed:\n  input:  %q\n  output: %q", line, got)
		}
	}
}

func TestParseDueDate(t *testing.T) {
	task := Parse("2026-02-23 請求書処理 due:2026-03-01")
	if task.DueDate.IsZero() {
		t.Fatal("expected DueDate to be set")
	}
	if !task.DueDate.Equal(date(2026, 3, 1)) {
		t.Errorf("expected due 2026-03-01, got %v", task.DueDate)
	}
	if task.Text != "請求書処理 due:2026-03-01" {
		t.Errorf("expected text to contain due:, got %q", task.Text)
	}
}

func TestParseDueDateWithTags(t *testing.T) {
	task := Parse("(A) 2026-02-23 企画書 +仕事 @PC due:2026-03-15")
	if !task.DueDate.Equal(date(2026, 3, 15)) {
		t.Errorf("expected due 2026-03-15, got %v", task.DueDate)
	}
	if task.Priority != "A" {
		t.Errorf("expected priority A, got %q", task.Priority)
	}
	if len(task.Projects) != 1 || task.Projects[0] != "仕事" {
		t.Errorf("expected project '仕事', got %v", task.Projects)
	}
}

func TestParseNoDueDate(t *testing.T) {
	task := Parse("2026-02-23 普通のタスク")
	if !task.DueDate.IsZero() {
		t.Errorf("expected no DueDate, got %v", task.DueDate)
	}
}

// --- DateTime tests ---

func datetime(y, m, d, h, min, sec int) time.Time {
	return time.Date(y, time.Month(m), d, h, min, sec, 0, time.UTC)
}

func TestParseWithDateTime(t *testing.T) {
	task := Parse("2026-02-23T14:30:00 レポートを書く")
	if !task.CreatedAt.Equal(datetime(2026, 2, 23, 14, 30, 0)) {
		t.Errorf("expected 2026-02-23T14:30:00, got %v", task.CreatedAt)
	}
	if !task.CreatedHasTime {
		t.Error("expected CreatedHasTime to be true")
	}
	if task.Text != "レポートを書く" {
		t.Errorf("expected 'レポートを書く', got %q", task.Text)
	}
}

func TestParseDateOnly_HasTimeFalse(t *testing.T) {
	task := Parse("2026-02-23 レポートを書く")
	if task.CreatedHasTime {
		t.Error("expected CreatedHasTime to be false for date-only")
	}
}

func TestParseDoneWithDateTime(t *testing.T) {
	task := Parse("x 2026-02-24T10:00:00 2026-02-23T14:30:00 牛乳を買う +買い物")
	if !task.Done {
		t.Error("expected done")
	}
	if !task.CompletedAt.Equal(datetime(2026, 2, 24, 10, 0, 0)) {
		t.Errorf("expected completed 2026-02-24T10:00:00, got %v", task.CompletedAt)
	}
	if !task.CompletedHasTime {
		t.Error("expected CompletedHasTime to be true")
	}
	if !task.CreatedAt.Equal(datetime(2026, 2, 23, 14, 30, 0)) {
		t.Errorf("expected created 2026-02-23T14:30:00, got %v", task.CreatedAt)
	}
	if !task.CreatedHasTime {
		t.Error("expected CreatedHasTime to be true")
	}
}

func TestSerializeWithDateTime(t *testing.T) {
	task := Task{
		CreatedAt:      datetime(2026, 2, 23, 14, 30, 0),
		CreatedHasTime: true,
		Text:           "レポートを書く",
	}
	got := Serialize(task)
	expected := "2026-02-23T14:30:00 レポートを書く"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSerializeDateOnly_NoTimeAppended(t *testing.T) {
	task := Task{
		CreatedAt:      date(2026, 2, 23),
		CreatedHasTime: false,
		Text:           "レポートを書く",
	}
	got := Serialize(task)
	expected := "2026-02-23 レポートを書く"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSerializeDoneWithDateTime(t *testing.T) {
	task := Task{
		Done:             true,
		CompletedAt:      datetime(2026, 2, 24, 10, 0, 0),
		CompletedHasTime: true,
		CreatedAt:        datetime(2026, 2, 23, 14, 30, 0),
		CreatedHasTime:   true,
		Text:             "牛乳を買う +買い物",
	}
	got := Serialize(task)
	expected := "x 2026-02-24T10:00:00 2026-02-23T14:30:00 牛乳を買う +買い物"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRoundTripDateTime(t *testing.T) {
	lines := []string{
		"2026-02-23T14:30:00 レポートを書く",
		"(A) 2026-02-23T09:00:00 企画書を書く +仕事 @PC",
		"x 2026-02-24T10:00:00 2026-02-23T14:30:00 牛乳を買う +買い物",
		"2026-02-23T14:30:00 請求書処理 +仕事 due:2026-03-01T18:00:00",
	}
	for _, line := range lines {
		task := Parse(line)
		got := Serialize(task)
		if got != line {
			t.Errorf("roundtrip failed:\n  input:  %q\n  output: %q", line, got)
		}
	}
}

// --- Note tests ---

func TestParseNote(t *testing.T) {
	task := Parse("x 2026-02-24 2026-02-23 牛乳を買う +買い物 note:やりきった！")
	if task.Note != "やりきった！" {
		t.Errorf("expected note 'やりきった！', got %q", task.Note)
	}
	if task.Text != "牛乳を買う +買い物" {
		t.Errorf("expected text '牛乳を買う +買い物', got %q", task.Text)
	}
}

func TestParseNoteWithSpaces(t *testing.T) {
	task := Parse("x 2026-02-24 2026-02-23 タスク note:達成感 ある！最高")
	if task.Note != "達成感 ある！最高" {
		t.Errorf("expected note '達成感 ある！最高', got %q", task.Note)
	}
}

func TestSerializeWithNote(t *testing.T) {
	task := Task{
		Done:        true,
		CompletedAt: date(2026, 2, 24),
		CreatedAt:   date(2026, 2, 23),
		Text:        "牛乳を買う +買い物",
		Note:        "やりきった！",
	}
	got := Serialize(task)
	expected := "x 2026-02-24 2026-02-23 牛乳を買う +買い物 note:やりきった！"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRoundTripWithNote(t *testing.T) {
	line := "x 2026-02-24T10:00:00 2026-02-23T14:30:00 牛乳を買う +買い物 note:やりきった！"
	task := Parse(line)
	got := Serialize(task)
	if got != line {
		t.Errorf("roundtrip failed:\n  input:  %q\n  output: %q", line, got)
	}
}

func TestParseNoNote(t *testing.T) {
	task := Parse("2026-02-23 普通のタスク")
	if task.Note != "" {
		t.Errorf("expected no note, got %q", task.Note)
	}
}

// --- DueDate with time tests ---

func TestParseDueDateWithTime(t *testing.T) {
	task := Parse("2026-02-23 プレゼン準備 due:2026-03-01T18:00:00")
	if !task.DueDate.Equal(datetime(2026, 3, 1, 18, 0, 0)) {
		t.Errorf("expected due 2026-03-01T18:00:00, got %v", task.DueDate)
	}
	if !task.DueHasTime {
		t.Error("expected DueHasTime to be true")
	}
}

func TestParseDueDateWithoutTime(t *testing.T) {
	task := Parse("2026-02-23 プレゼン準備 due:2026-03-01")
	if !task.DueDate.Equal(date(2026, 3, 1)) {
		t.Errorf("expected due 2026-03-01, got %v", task.DueDate)
	}
	if task.DueHasTime {
		t.Error("expected DueHasTime to be false")
	}
}
