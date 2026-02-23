package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hengin-eer/todome/internal/todo"
	"github.com/hengin-eer/todome/internal/ui"
	"github.com/spf13/cobra"
)

var (
	listAll     bool
	listDone    bool
	listUndone  bool
	listOverdue bool
	listOr      bool
	listNot     bool
	listSort    string
	listReverse bool
)

var listCmd = &cobra.Command{
	Use:   "list [+project...] [@context...]",
	Short: "タスク一覧を表示する",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate exclusive flags
		flagCount := 0
		if listAll {
			flagCount++
		}
		if listDone {
			flagCount++
		}
		if listUndone {
			flagCount++
		}
		if flagCount > 1 {
			return fmt.Errorf("--all, --done, --undone は同時に指定できない")
		}

		s := getStore()
		tasks, err := s.Load()
		if err != nil {
			return fmt.Errorf("読み込みエラー: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("タスクがないぞ。todome add で追加しろ！")
			return nil
		}

		// Parse filter args
		var filterProjects, filterContexts []string
		for _, arg := range args {
			if strings.HasPrefix(arg, "+") {
				filterProjects = append(filterProjects, arg[1:])
			} else if strings.HasPrefix(arg, "@") {
				filterContexts = append(filterContexts, arg[1:])
			}
		}

		// Build indexed list for stable numbering
		var filtered []indexedTask

		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		for i, t := range tasks {
			// State filter
			switch {
			case listDone && !t.Done:
				continue
			case listUndone && t.Done:
				continue
			case listOverdue:
				if t.Done || t.DueDate.IsZero() || !t.DueDate.Before(today) {
					continue
				}
			case !listAll && !listDone && t.Done:
				continue
			}

			// Tag filter
			if len(filterProjects) > 0 || len(filterContexts) > 0 {
				if !matchTags(t, filterProjects, filterContexts, listOr, listNot) {
					continue
				}
			}

			filtered = append(filtered, indexedTask{num: i + 1, task: t})
		}

		// Sort
		if listSort != "" {
			sortTasks(filtered, listSort, listReverse)
		}

		if len(filtered) == 0 {
			fmt.Println("条件に一致するタスクがない")
			return nil
		}

		for _, it := range filtered {
			fmt.Println(ui.FormatTask(it.num, it.task))
		}
		fmt.Printf("\n%d 件のタスク\n", len(filtered))
		return nil
	},
}

func matchTags(t todo.Task, projects, contexts []string, orMode, notMode bool) bool {
	matched := tagMatch(t, projects, contexts, orMode)
	if notMode {
		return !matched
	}
	return matched
}

func tagMatch(t todo.Task, projects, contexts []string, orMode bool) bool {
	if orMode {
		for _, p := range projects {
			if containsStr(t.Projects, p) {
				return true
			}
		}
		for _, c := range contexts {
			if containsStr(t.Contexts, c) {
				return true
			}
		}
		return false
	}
	// AND mode
	for _, p := range projects {
		if !containsStr(t.Projects, p) {
			return false
		}
	}
	for _, c := range contexts {
		if !containsStr(t.Contexts, c) {
			return false
		}
	}
	return true
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func sortTasks(tasks []indexedTask, field string, reverse bool) {
	sort.SliceStable(tasks, func(i, j int) bool {
		var less bool
		switch field {
		case "priority":
			less = comparePriority(tasks[i].task.Priority, tasks[j].task.Priority)
		case "created":
			less = tasks[i].task.CreatedAt.After(tasks[j].task.CreatedAt)
		case "due":
			less = compareDue(tasks[i].task.DueDate, tasks[j].task.DueDate)
		default:
			return false
		}
		if reverse {
			return !less
		}
		return less
	})
}

// comparePriority: A < B < C < ... < "" (no priority last)
func comparePriority(a, b string) bool {
	if a == b {
		return false
	}
	if a == "" {
		return false
	}
	if b == "" {
		return true
	}
	return a < b
}

// compareDue: earlier due first, no due date last
func compareDue(a, b time.Time) bool {
	if a.IsZero() && b.IsZero() {
		return false
	}
	if a.IsZero() {
		return false
	}
	if b.IsZero() {
		return true
	}
	return a.Before(b)
}

type indexedTask struct {
	num  int
	task todo.Task
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "完了タスクも表示する")
	listCmd.Flags().BoolVar(&listDone, "done", false, "完了タスクのみ表示する")
	listCmd.Flags().BoolVar(&listUndone, "undone", false, "未完了タスクのみ表示する")
	listCmd.Flags().BoolVar(&listOverdue, "overdue", false, "期限切れタスクのみ表示する")
	listCmd.Flags().BoolVar(&listOr, "or", false, "フィルタをOR結合にする")
	listCmd.Flags().BoolVarP(&listNot, "not", "n", false, "フィルタを除外条件にする")
	listCmd.Flags().StringVarP(&listSort, "sort", "s", "", "ソート: priority, created, due")
	listCmd.Flags().BoolVarP(&listReverse, "reverse", "r", false, "ソート順を反転する")
	rootCmd.AddCommand(listCmd)
}
