package transport

type AliceRequest struct {
	Meta    AliceMeta           `json:"meta"`
	Session AliceRequestSession `json:"session"`
	Request AliceUtterance      `json:"request"`
	Version string              `json:"version"`
}

type AliceMeta struct {
	Locale     string          `json:"locale"`
	Timezone   string          `json:"timezone"`
	ClientID   string          `json:"client_id"`
	Interfaces AliceInterfaces `json:"interfaces"`
}

type AliceInterfaces struct {
	Screen         map[string]interface{} `json:"screen"`
	Payments       map[string]interface{} `json:"payments"`
	AccountLinking map[string]interface{} `json:"account_linking"`
}

type AliceRequestSession struct {
	MessageID   int              `json:"message_id"`
	SessionID   string           `json:"session_id"`
	SkillID     string           `json:"skill_id"`
	User        AliceUser        `json:"user"`
	Application AliceApplication `json:"application"`
	UserID      string           `json:"user_id"`
	New         bool             `json:"new"`
}

type AliceUser struct {
	UserID string `json:"user_id"`
}

type AliceApplication struct {
	ApplicationID string `json:"application_id"`
}

type AliceUtterance struct {
	Command           string      `json:"command"`
	OriginalUtterance string      `json:"original_utterance"`
	NLU               AliceNLU    `json:"nlu"`
	Markup            AliceMarkup `json:"markup"`
	Type              string      `json:"type"`
}

type AliceNLU struct {
	Tokens   []string                 `json:"tokens"`
	Entities []map[string]interface{} `json:"entities"`
	Intents  map[string]interface{}   `json:"intents"`
}

type AliceMarkup struct {
	DangerousContext bool `json:"dangerous_context"`
}

type AliceResponse struct {
	Response AliceResponseBody    `json:"response"`
	Session  AliceResponseSession `json:"session"`
	Version  string               `json:"version"`
}

type AliceResponseBody struct {
	Text       string `json:"text"`
	EndSession bool   `json:"end_session"`
}

type AliceResponseSession struct {
	SessionID string `json:"session_id"`
	MessageID int    `json:"message_id"`
}
