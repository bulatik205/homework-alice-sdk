package alice

import (
	"context"
	"strings"
	"time"

	homework "github.com/bulatik205/homework-sdk-go"
)

func handleHomework(ctx context.Context, cfg Config, req Request, cmd string) Response {
	loc := loadLocation(req.Meta.Timezone)

	if query := extractSubjectQuery(cmd); query != "" {
		return handleBySubject(ctx, cfg, query)
	}

	date := resolveDate(cmd, loc)
	if !hasExplicitDate(cmd) {
		date = tomorrow(loc)
	}

	tasks, err := cfg.HW.TasksByDate(ctx, date)
	if err != nil {
		return simpleResp(apiErrorText)
	}

	label := dateLabel(cmd, date, loc)
	return renderTasks(ctx, cfg, tasks, label, loc)
}

func handleBySubject(ctx context.Context, cfg Config, query string) Response {
	subjects, err := cfg.HW.Subjects(ctx)
	if err != nil {
		return simpleResp(apiErrorText)
	}

	matched := matchSubject(subjects, query)
	if matched == nil {
		return simpleResp(subjectNotFound)
	}

	task, err := cfg.HW.TaskBySubject(ctx, matched.Code)
	if err != nil {
		return simpleResp(apiErrorText)
	}
	if task == nil {
		return simpleResp("По предмету " + matched.Display + " ничего не задали.")
	}

	return renderTasks(ctx, cfg, []homework.Task{*task}, "по предмету "+matched.Display, nil)
}

func handleButton(ctx context.Context, cfg Config, req Request) Response {
	loc := loadLocation(req.Meta.Timezone)

	if d, ok := req.Request.Payload["date"].(string); ok {
		tasks, err := cfg.HW.TasksByDate(ctx, d)
		if err != nil {
			return simpleResp(apiErrorText)
		}
		return renderTasks(ctx, cfg, tasks, d, loc)
	}
	if s, ok := req.Request.Payload["subject"].(string); ok {
		task, err := cfg.HW.TaskBySubject(ctx, s)
		if err != nil {
			return simpleResp(apiErrorText)
		}
		if task == nil {
			return simpleResp("По этому предмету ничего не задали.")
		}
		return renderTasks(ctx, cfg, []homework.Task{*task}, "по предмету "+s, nil)
	}
	return simpleResp("Не поняла, что показать.")
}

func renderTasks(ctx context.Context, cfg Config, tasks []homework.Task, label string, loc *time.Location) Response {
	if len(tasks) == 0 {
		return simpleResp("На " + label + " ничего не задали.")
	}

	display := subjectDisplays(ctx, cfg)

	items := make([]CardItem, 0, len(tasks))
	for _, t := range tasks {
		subj := display[t.Subject]
		if subj == "" {
			subj = t.Subject
		}
		items = append(items, CardItem{
			Title:       subj,
			Description: t.Task,
		})
	}
	if len(items) > 5 {
		items = items[:5]
	}

	n := cfg.MaxVoiceItems
	if n > len(tasks) {
		n = len(tasks)
	}
	var sb strings.Builder
	sb.WriteString("На " + label + ": " + plural(tasksCount(len(tasks))) + ". ")
	for i := 0; i < n; i++ {
		subj := display[tasks[i].Subject]
		if subj == "" {
			subj = tasks[i].Subject
		}
		sb.WriteString(subj + ": " + tasks[i].Task + ". ")
	}
	if len(tasks) > n {
		sb.WriteString("Остальное на экране.")
	}

	var buttons []Button
	if loc != nil {
		buttons = []Button{
			{Title: "На сегодня", Payload: map[string]any{"date": today(loc)}},
			{Title: "На завтра", Payload: map[string]any{"date": tomorrow(loc)}},
		}
	}

	return itemsListResp(sb.String(), "Домашка", items, buttons)
}

func subjectDisplays(ctx context.Context, cfg Config) map[string]string {
	subjects, err := cfg.HW.Subjects(ctx)
	if err != nil {
		return nil
	}
	m := make(map[string]string, len(subjects))
	for _, s := range subjects {
		m[s.Code] = s.Display
	}
	return m
}

func matchSubject(subjects []homework.Subject, query string) *homework.Subject {
	q := strings.ToLower(query)
	for i := range subjects {
		d := strings.ToLower(subjects[i].Display)
		if strings.Contains(d, q) || strings.Contains(q, d) {
			return &subjects[i]
		}
	}
	return nil
}

func extractSubjectQuery(cmd string) string {
	idx := strings.Index(cmd, " по ")
	if idx < 0 {
		return ""
	}
	rest := strings.TrimSpace(cmd[idx+len(" по "):])
	for _, stop := range []string{" на ", " за ", " к "} {
		if i := strings.Index(rest, stop); i > 0 {
			rest = strings.TrimSpace(rest[:i])
		}
	}
	return rest
}

func loadLocation(tz string) *time.Location {
	if tz == "" {
		tz = "Europe/Moscow"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc, _ = time.LoadLocation("Europe/Moscow")
	}
	return loc
}

func dateLabel(cmd, date string, loc *time.Location) string {
	now := time.Now().In(loc)
	switch date {
	case now.Format("2006-01-02"):
		return "сегодня"
	case now.AddDate(0, 0, 1).Format("2006-01-02"):
		return "завтра"
	case now.AddDate(0, 0, 2).Format("2006-01-02"):
		return "послезавтра"
	}
	return date
}

func today(loc *time.Location) string {
	return time.Now().In(loc).Format("2006-01-02")
}

func tomorrow(loc *time.Location) string {
	return time.Now().In(loc).AddDate(0, 0, 1).Format("2006-01-02")
}

func tasksCount(n int) int { return n }

func plural(n int) string {
	mod10 := n % 10
	mod100 := n % 100
	switch {
	case mod10 == 1 && mod100 != 11:
		return "1 предмет"
	case mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20):
		return itoa(n) + " предмета"
	default:
		return itoa(n) + " предметов"
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
