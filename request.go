package alice

type Request struct {
	Meta    Meta    `json:"meta"`
	Request Req     `json:"request"`
	Session Session `json:"session"`
	State   State   `json:"state,omitempty"`
	Version string  `json:"version"`
}

type Meta struct {
	Locale     string     `json:"locale"`
	Timezone   string     `json:"timezone"`
	ClientID   string     `json:"client_id"`
	Interfaces Interfaces `json:"interfaces"`
}

type Interfaces struct {
	Screen         *struct{} `json:"screen,omitempty"`
	AccountLinking *struct{} `json:"account_linking,omitempty"`
	AudioPlayer    *struct{} `json:"audio_player,omitempty"`
}

type Session struct {
	MessageID   int          `json:"message_id"`
	SessionID   string       `json:"session_id"`
	SkillID     string       `json:"skill_id"`
	UserID      string       `json:"user_id"`
	User        *SessionUser `json:"user,omitempty"`
	Application Application  `json:"application"`
	New         bool         `json:"new"`
}

type SessionUser struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token,omitempty"`
}

type Application struct {
	ApplicationID string `json:"application_id"`
}

type State struct {
	Session     map[string]any `json:"session,omitempty"`
	User        map[string]any `json:"user,omitempty"`
	Application map[string]any `json:"application,omitempty"`
}

type Req struct {
	Type              string         `json:"type"`
	Command           string         `json:"command,omitempty"`
	OriginalUtterance string         `json:"original_utterance,omitempty"`
	Payload           map[string]any `json:"payload,omitempty"`
	NLU               *NLU           `json:"nlu,omitempty"`
}

type NLU struct {
	Tokens   []string          `json:"tokens"`
	Entities []Entity          `json:"entities"`
	Intents  map[string]Intent `json:"intents"`
}

type Entity struct {
	Tokens []string `json:"tokens"`
	Type   string   `json:"type"`
	Value  any      `json:"value"`
}

type Intent struct {
	Slots map[string]Slot `json:"slots"`
}

type Slot struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}
