package alice

import (
	"context"
	"strings"
)

func route(ctx context.Context, cfg Config, req Request) Response {
	if req.Request.Type == "ButtonPressed" && req.Request.Payload != nil {
		return handleButton(ctx, cfg, req)
	}

	cmd := strings.ToLower(strings.TrimSpace(req.Request.Command))

	switch {
	case cmd == "":
		return simpleResp("Не расслышала. Повтори, пожалуйста.")

	case containsAny(cmd, "привет", "здравствуй", "хай", "добрый день", "добрый вечер"):
		return simpleResp("Привет! Спроси, что задали: «домашка», «домашка по математике» или «домашка на завтра».")

	case containsAny(cmd, "помощь", "что умеешь", "как пользоваться", "что ты можешь"):
		return helpResp()

	case containsAny(cmd, "домашк", "дз", "задали", "задание", "уроки"):
		return handleHomework(ctx, cfg, req, cmd)

	default:
		return simpleResp("Не поняла. Скажи «домашка» или «домашка по математике».")
	}
}

func helpResp() Response {
	return simpleResp(
		"Я умею рассказывать домашку. Скажи: «домашка» — на сегодня, " +
			"«домашка на завтра» — на завтра, или «домашка по математике» — по предмету.",
	)
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
