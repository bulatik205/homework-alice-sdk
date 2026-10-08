package alice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	homework "github.com/bulatik205/homework-sdk-go"
)

type Config struct {
	HW            *homework.Client
	MaxVoiceItems int
}

func Run(cfg Config) http.HandlerFunc {
	if cfg.MaxVoiceItems == 0 {
		cfg.MaxVoiceItems = 3
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, simpleResp("Не разобрала запрос, попробуй ещё раз."))
			return
		}

		if req.Request.OriginalUtterance == "" && req.Request.Command == "" {
			writeJSON(w, simpleResp("Привет! Спроси, что задали на завтра."))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		resp := route(ctx, cfg, req)
		writeJSON(w, resp)
	}
}

func writeJSON(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(resp)
}
