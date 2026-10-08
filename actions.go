package alice

import (
	"context"
	"strings"
	"time"

	homework "github.com/bulatik205/homework-sdk-go"
)

var subjectAliases = map[string][]string{
	"math":            {"математика", "матеша", "матеш", "матан", "алгебра"},
	"geometry":        {"геометрия", "геом", "геома"},
	"russian":         {"русский", "русский язык", "русс"},
	"literature":      {"литература", "литра", "лит"},
	"history":         {"история", "ист"},
	"physics":         {"физика", "физра", "физ"},
	"chemistry":       {"химия", "хим"},
	"biology":         {"биология", "био"},
	"geography":       {"география", "гео"},
	"english":         {"английский", "английский язык", "англ", "инглиш"},
	"social":          {"обществознание", "общество"},
	"random-and-stat": {"вероятность и статистика", "вероятность", "статистика"},
	"projects":        {"проектная деятельность", "проект", "проекты"},
}

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
	return renderTasks(ctx, cfg, tasks, label, loc, "date")
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

	loc := time.Now().Location()
	return renderTasks(ctx, cfg, []homework.Task{*task}, matched.Display, loc, "subject")
}

func handleButton(ctx context.Context, cfg Config, req Request) Response {
	loc := loadLocation(req.Meta.Timezone)

	if d, ok := req.Request.Payload["date"].(string); ok {
		tasks, err := cfg.HW.TasksByDate(ctx, d)
		if err != nil {
			return simpleResp(apiErrorText)
		}
		return renderTasks(ctx, cfg, tasks, d, loc, "date")
	}
	if s, ok := req.Request.Payload["subject"].(string); ok {
		task, err := cfg.HW.TaskBySubject(ctx, s)
		if err != nil {
			return simpleResp(apiErrorText)
		}
		if task == nil {
			return simpleResp("По этому предмету ничего не задали.")
		}
		return renderTasks(ctx, cfg, []homework.Task{*task}, s, loc, "subject")
	}
	return simpleResp("Не поняла, что показать.")
}

func renderTasks(ctx context.Context, cfg Config, tasks []homework.Task, label string, loc *time.Location, kind string) Response {
	if len(tasks) == 0 {
		if kind == "subject" {
			return simpleResp("По предмету " + label + " ничего не задали.")
		}
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
	if kind == "subject" {
		sb.WriteString("По предмету " + label + ": " + plural(len(tasks)) + ". ")
	} else {
		sb.WriteString("На " + label + ": " + plural(len(tasks)) + ". ")
	}

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
	if loc != nil && kind == "date" {
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
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}

	for i := range subjects {
		d := strings.ToLower(subjects[i].Display)
		if matchByStem(d, q) {
			return &subjects[i]
		}
	}

	for i := range subjects {
		for _, alias := range subjectAliases[subjects[i].Code] {
			if matchByStem(strings.ToLower(alias), q) {
				return &subjects[i]
			}
		}
	}

	return nil
}

func matchByStem(a, b string) bool {
	const minLen = 4
	if len(a) < minLen || len(b) < minLen {
		return a == b
	}
	return a[:minLen] == b[:minLen]
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
