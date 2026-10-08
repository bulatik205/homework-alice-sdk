package alice

import "fmt"

const (
	emptyTasksText    = "Я не нашла домашнего задания."
	apiErrorText      = "Сервис домашки сейчас не отвечает. Попробуй чуть позже."
	subjectNotFound   = "Не нашла такой предмет. Попробуй назвать по-другому."
	noSubjectForQuery = "По какому предмету? Например, «домашка по математике»."
)

func tasksHeader(date string, n int) string {
	return fmt.Sprintf("На %s: %d.", date, n)
}
