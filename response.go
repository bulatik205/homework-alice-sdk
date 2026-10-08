package alice

type Response struct {
	Response         ResponseBody   `json:"response"`
	SessionState     map[string]any `json:"session_state,omitempty"`
	UserStateUpdate  map[string]any `json:"user_state_update,omitempty"`
	ApplicationState map[string]any `json:"application_state,omitempty"`
	Version          string         `json:"version"`
}

type ResponseBody struct {
	Text       string      `json:"text"`
	TTS        string      `json:"tts,omitempty"`
	Card       *Card       `json:"card,omitempty"`
	Buttons    []Button    `json:"buttons,omitempty"`
	EndSession bool        `json:"end_session"`
	Directives *Directives `json:"directives,omitempty"`
}

type Directives struct {
	StartAccountLinking *struct{} `json:"start_account_linking,omitempty"`
}

type Card struct {
	Type        string     `json:"type"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	ImageID     string     `json:"image_id,omitempty"`
	Items       []CardItem `json:"items,omitempty"`
}

type CardItem struct {
	ImageID     string  `json:"image_id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Button      *Button `json:"button,omitempty"`
}

type Button struct {
	Title   string         `json:"title"`
	Payload map[string]any `json:"payload,omitempty"`
	URL     string         `json:"url,omitempty"`
	Hide    bool           `json:"hide,omitempty"`
}

func simpleResp(text string) Response {
	return Response{
		Response: ResponseBody{Text: text, EndSession: false},
		Version:  "1.0",
	}
}

func itemsListResp(text, title string, items []CardItem, buttons []Button) Response {
	return Response{
		Response: ResponseBody{
			Text:       text,
			Card:       &Card{Type: "ItemsList", Title: title, Items: items},
			Buttons:    buttons,
			EndSession: false,
		},
		Version: "1.0",
	}
}
