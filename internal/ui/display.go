package ui

import (
	"fmt"
	"math"
	"time"

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

	bgRedWhite    = "\033[41;37;1m" // 赤背景+白文字+太字
	bgYellowBlack = "\033[43;30m"   // 黄背景+黒文字
)

// FormatTask formats a task for display with line number.
func FormatTask(num int, t todo.Task) string {
	prefix := fmt.Sprintf("%s%2d.%s ", bold, num, reset)

	if t.Done {
		line := gray + "✓ " + todo.Serialize(t) + reset
		if t.Note != "" {
			line += "\n     " + gray + "📝 " + t.Note + reset
		}
		return prefix + line
	}

	line := todo.Serialize(t)
	if t.Priority != "" {
		color := priorityColor(t.Priority)
		line = color + line + reset
	}

	if suffix := dueSuffix(t); suffix != "" {
		line += " " + suffix
	}

	return prefix + line
}

func dueSuffix(t todo.Task) string {
	if t.DueDate.IsZero() || t.Done {
		return ""
	}
	now := time.Now()

	if t.DueHasTime {
		diff := time.Until(t.DueDate)
		switch {
		case diff < 0:
			return bgRedWhite + " [期限切れ!] " + reset
		case diff < time.Hour:
			mins := int(diff.Minutes())
			return bgYellowBlack + fmt.Sprintf(" [あと%d分] ", mins) + reset
		case diff < 24*time.Hour:
			hours := int(diff.Hours())
			return bgYellowBlack + fmt.Sprintf(" [あと%d時間] ", hours) + reset
		case diff < 3*24*time.Hour:
			days := int(diff.Hours() / 24)
			return bgYellowBlack + fmt.Sprintf(" [あと%d日] ", days) + reset
		default:
			return ""
		}
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	due := time.Date(t.DueDate.Year(), t.DueDate.Month(), t.DueDate.Day(), 0, 0, 0, 0, time.Local)
	days := int(math.Ceil(due.Sub(today).Hours() / 24))

	switch {
	case days < 0:
		return bgRedWhite + " [期限切れ!] " + reset
	case days == 0:
		return bgYellowBlack + " [今日まで] " + reset
	case days <= 3:
		return bgYellowBlack + fmt.Sprintf(" [あと%d日] ", days) + reset
	default:
		return ""
	}
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
