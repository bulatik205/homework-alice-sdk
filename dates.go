package alice

import (
	"strconv"
	"strings"
	"time"
)

func resolveDate(cmd string, loc *time.Location) string {
	now := time.Now().In(loc)

	switch {
	case containsAny(cmd, "послезавтра"):
		return now.AddDate(0, 0, 2).Format("2006-01-02")
	case containsAny(cmd, "завтра"):
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	case containsAny(cmd, "вчера"):
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	case containsAny(cmd, "сегодня"):
		return now.Format("2006-01-02")
	}

	if d, ok := parseExplicitDate(cmd, now); ok {
		return d
	}
	return now.Format("2006-01-02")
}

func parseExplicitDate(cmd string, now time.Time) (string, bool) {
	for _, token := range strings.Fields(cmd) {
		if t, err := time.Parse("2006-01-02", token); err == nil {
			return t.Format("2006-01-02"), true
		}
		if t, err := time.Parse("02.01.2006", token); err == nil {
			return t.Format("2006-01-02"), true
		}
	}

	months := map[string]time.Month{
		"январ": time.January, "феврал": time.February, "март": time.March,
		"апрел": time.April, "ма": time.May, "июн": time.June,
		"июл": time.July, "август": time.August, "сентябр": time.September,
		"октябр": time.October, "ноябр": time.November, "декабр": time.December,
	}
	for prefix, month := range months {
		idx := strings.Index(cmd, prefix)
		if idx <= 0 {
			continue
		}
		before := strings.Fields(cmd[:idx])
		if len(before) == 0 {
			continue
		}
		day, err := strconv.Atoi(before[len(before)-1])
		if err != nil || day < 1 || day > 31 {
			continue
		}
		year := now.Year()
		if month < now.Month() {
			year++
		}
		return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Format("2006-01-02"), true
	}
	return "", false
}
