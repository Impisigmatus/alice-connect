// Этот тестовый файл сгенерирован автоматически.
package models

import (
	"encoding/json"
	"testing"
)

func TestAliceResponseJSON(t *testing.T) {
	response := AliceResponse{
		Response: AliceResponseBody{Text: "Привет", EndSession: true},
		Session:  AliceResponseSession{SessionID: "test-session", MessageID: 1},
		Version:  "1.0",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["version"] != "1.0" || decoded["response"] == nil || decoded["session"] == nil {
		t.Fatalf("неверный JSON-ответ: %s", data)
	}
}
