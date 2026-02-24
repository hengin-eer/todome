package todo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	dateFormat     = "2006-01-02"
	dateTimeFormat = "2006-01-02T15:04:05"
)

var (
	projectRe = regexp.MustCompile(`(?:^|\s)\+(\S+)`)
	contextRe = regexp.MustCompile(`(?:^|\s)@(\S+)`)
	dueRe     = regexp.MustCompile(`(?:^|\s)due:(\d{4}-\d{2}-\d{2}(?:T\d{2}:\d{2}:\d{2})?)(?:\s|$)`)
	noteRe    = regexp.MustCompile(`(?:^|\s)note:(.+)$`)
)

// Parse parses a single todo.txt line into a Task.
func Parse(line string) Task {
	t := Task{}
	s := line

	// 1. Check done
	if strings.HasPrefix(s, "x ") {
		t.Done = true
		s = s[2:]
	}

	// 2. Priority (only for non-done tasks in standard format)
	if !t.Done && len(s) >= 4 && s[0] == '(' && s[2] == ')' && s[1] >= 'A' && s[1] <= 'Z' && s[3] == ' ' {
		t.Priority = string(s[1])
		s = s[4:]
	}

	// 3. Dates (try datetime first, then date)
	if t.Done {
		if d, rest, hasTime, ok := parseLeadingDateTime(s); ok {
			t.CompletedAt = d
			t.CompletedHasTime = hasTime
			s = rest
			if d2, rest2, hasTime2, ok := parseLeadingDateTime(s); ok {
				t.CreatedAt = d2
				t.CreatedHasTime = hasTime2
				s = rest2
			}
		}
	} else {
		if d, rest, hasTime, ok := parseLeadingDateTime(s); ok {
			t.CreatedAt = d
			t.CreatedHasTime = hasTime
			s = rest
		}
	}

	// 4. Extract note (before tag extraction to avoid false +/@ matches in notes)
	if m := noteRe.FindStringSubmatch(s); m != nil {
		t.Note = m[1]
		s = strings.TrimSpace(noteRe.ReplaceAllString(s, ""))
	}

	// 5. Extract +project and @context tags, and due date
	for _, m := range projectRe.FindAllStringSubmatch(s, -1) {
		t.Projects = append(t.Projects, m[1])
	}
	for _, m := range contextRe.FindAllStringSubmatch(s, -1) {
		t.Contexts = append(t.Contexts, m[1])
	}
	if m := dueRe.FindStringSubmatch(s); m != nil {
		dueStr := m[1]
		if len(dueStr) > 10 {
			if d, err := time.Parse(dateTimeFormat, dueStr); err == nil {
				t.DueDate = d
				t.DueHasTime = true
			}
		} else {
			if d, err := time.Parse(dateFormat, dueStr); err == nil {
				t.DueDate = d
			}
		}
	}

	t.Text = strings.TrimSpace(s)
	return t
}

// Serialize converts a Task back to a todo.txt line.
func Serialize(t Task) string {
	var parts []string

	if t.Done {
		parts = append(parts, "x")
		if !t.CompletedAt.IsZero() {
			if t.CompletedHasTime {
				parts = append(parts, t.CompletedAt.Format(dateTimeFormat))
			} else {
				parts = append(parts, t.CompletedAt.Format(dateFormat))
			}
		}
		if !t.CreatedAt.IsZero() {
			if t.CreatedHasTime {
				parts = append(parts, t.CreatedAt.Format(dateTimeFormat))
			} else {
				parts = append(parts, t.CreatedAt.Format(dateFormat))
			}
		}
	} else {
		if t.Priority != "" {
			parts = append(parts, fmt.Sprintf("(%s)", t.Priority))
		}
		if !t.CreatedAt.IsZero() {
			if t.CreatedHasTime {
				parts = append(parts, t.CreatedAt.Format(dateTimeFormat))
			} else {
				parts = append(parts, t.CreatedAt.Format(dateFormat))
			}
		}
	}

	parts = append(parts, t.Text)

	if t.Note != "" {
		parts = append(parts, "note:"+t.Note)
	}

	return strings.Join(parts, " ")
}

// parseLeadingDateTime tries to parse a leading datetime (YYYY-MM-DDTHH:mm:ss) or date (YYYY-MM-DD).
func parseLeadingDateTime(s string) (time.Time, string, bool, bool) {
	// Try datetime first (19 chars: 2006-01-02T15:04:05)
	if len(s) >= 19 {
		if d, err := time.Parse(dateTimeFormat, s[:19]); err == nil {
			rest := s[19:]
			if len(rest) > 0 && rest[0] == ' ' {
				rest = rest[1:]
			}
			return d, rest, true, true
		}
	}
	// Fall back to date only (10 chars: 2006-01-02)
	if len(s) >= 10 {
		if d, err := time.Parse(dateFormat, s[:10]); err == nil {
			rest := s[10:]
			if len(rest) > 0 && rest[0] == ' ' {
				rest = rest[1:]
			}
			return d, rest, false, true
		}
	}
	return time.Time{}, s, false, false
}
