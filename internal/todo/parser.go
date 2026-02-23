package todo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

var (
	projectRe = regexp.MustCompile(`(?:^|\s)\+(\S+)`)
	contextRe = regexp.MustCompile(`(?:^|\s)@(\S+)`)
	dueRe     = regexp.MustCompile(`(?:^|\s)due:(\d{4}-\d{2}-\d{2})(?:\s|$)`)
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

	// 3. Dates
	if t.Done {
		// completed: completion_date creation_date
		if d, rest, ok := parseLeadingDate(s); ok {
			t.CompletedAt = d
			s = rest
			if d2, rest2, ok := parseLeadingDate(s); ok {
				t.CreatedAt = d2
				s = rest2
			}
		}
	} else {
		if d, rest, ok := parseLeadingDate(s); ok {
			t.CreatedAt = d
			s = rest
		}
	}

	// 4. Extract +project and @context tags, and due date
	for _, m := range projectRe.FindAllStringSubmatch(s, -1) {
		t.Projects = append(t.Projects, m[1])
	}
	for _, m := range contextRe.FindAllStringSubmatch(s, -1) {
		t.Contexts = append(t.Contexts, m[1])
	}
	if m := dueRe.FindStringSubmatch(s); m != nil {
		if d, err := time.Parse(dateFormat, m[1]); err == nil {
			t.DueDate = d
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
			parts = append(parts, t.CompletedAt.Format(dateFormat))
		}
		if !t.CreatedAt.IsZero() {
			parts = append(parts, t.CreatedAt.Format(dateFormat))
		}
	} else {
		if t.Priority != "" {
			parts = append(parts, fmt.Sprintf("(%s)", t.Priority))
		}
		if !t.CreatedAt.IsZero() {
			parts = append(parts, t.CreatedAt.Format(dateFormat))
		}
	}

	parts = append(parts, t.Text)
	return strings.Join(parts, " ")
}

func parseLeadingDate(s string) (time.Time, string, bool) {
	if len(s) < 10 {
		return time.Time{}, s, false
	}
	d, err := time.Parse(dateFormat, s[:10])
	if err != nil {
		return time.Time{}, s, false
	}
	rest := s[10:]
	if len(rest) > 0 && rest[0] == ' ' {
		rest = rest[1:]
	}
	return d, rest, true
}
