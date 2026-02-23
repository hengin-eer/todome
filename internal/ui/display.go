package ui

import (
	"fmt"

	"github.com/hengin-eer/todome/internal/todo"
)

// ANSI color codes
const (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	bold   = "\033[1m"
)

// FormatTask formats a task for display with line number.
func FormatTask(num int, t todo.Task) string {
	prefix := fmt.Sprintf("%s%2d.%s ", bold, num, reset)

	if t.Done {
		return prefix + gray + "✓ " + todo.Serialize(t) + reset
	}

	if t.Priority != "" {
		color := priorityColor(t.Priority)
		return prefix + color + todo.Serialize(t) + reset
	}

	return prefix + todo.Serialize(t)
}

func priorityColor(p string) string {
	switch p {
	case "A":
		return red
	case "B":
		return yellow
	case "C":
		return green
	default:
		return ""
	}
}
