// Этот тестовый файл сгенерирован автоматически.
package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Impisigmatus/alice-connect/internal/config"
	"github.com/Impisigmatus/alice-connect/internal/models"
	"github.com/Impisigmatus/alice-connect/internal/usecase"
	"github.com/Impisigmatus/service_core/log"
	"github.com/rs/zerolog"
)

func TestFormatSystemStatus(t *testing.T) {
	status := &models.SystemStatus{
		CPUPercent: 12.5, MemoryPercent: 25, MemoryUsed: 512 << 20, MemoryTotal: 2048 << 20,
		Temperatures: []models.SystemTemperature{{Sensor: "cpu_thermal", Celsius: 47.4}},
		SerialData:   "Привет из Arduino! Время работы: 42 с",
	}
	text := formatSystemStatus(status)
	for _, want := range []string{"ЦПУ: 12.5%", "ОЗУ: 25.0% (512 из 2048 МБ)", "Температура cpu_thermal: 47°C", "Arduino: Привет из Arduino! Время работы: 42 с"} {
		if !strings.Contains(text, want) {
			t.Errorf("в ответе %q нет %q", text, want)
		}
	}
}

func TestPostApiAliceWebhookRejectsInvalidRequest(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: "{"},
		{name: "missing session ID", body: `{"session":{}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/alice/webhook", strings.NewReader(tc.body))
			req = req.WithContext(context.WithValue(req.Context(), log.CtxKey, zerolog.Nop()))
			response := httptest.NewRecorder()
			(&Transport{}).PostApiAliceWebhook(response, req)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("HTTP %d; ожидался 400", response.Code)
			}
		})
	}
}

func TestPostApiAliceWebhookResponse(t *testing.T) {
	serialServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Arduino работает"))
	}))
	defer serialServer.Close()
	uc := usecase.New(config.SerialDriverAPIConfig{Hostname: serialServer.URL, Timeout: 2 * time.Second})

	for _, tc := range []struct {
		name    string
		command string
		want    string
	}{
		{name: "запуск навыка", want: "Привет! Go сервис на Orange Pi работает."},
		{name: "загрузка процессора", command: "какая загрузка процессора", want: "ЦПУ:"},
		{name: "короткая команда", command: "загрузка процессора", want: "ЦПУ:"},
		{name: "регистр и пробелы", command: "  Какая   ЗАГРУЗКА\tпроцессора  ", want: "ЦПУ:"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, err := json.Marshal(map[string]any{
				"request": map[string]string{"command": tc.command},
				"session": map[string]any{"session_id": "test-session", "message_id": 3},
			})
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/alice/webhook", strings.NewReader(string(requestBody)))
			req = req.WithContext(context.WithValue(req.Context(), log.CtxKey, zerolog.Nop()))
			response := httptest.NewRecorder()
			New("1.0", uc).PostApiAliceWebhook(response, req)
			if response.Code != http.StatusOK {
				t.Fatalf("HTTP %d: %s", response.Code, response.Body.String())
			}
			var body models.AliceResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Version != "1.0" || body.Session.SessionID != "test-session" || body.Session.MessageID != 3 || !body.Response.EndSession {
				t.Fatalf("неверный ответ Алисе: %+v", body)
			}
			if !strings.Contains(body.Response.Text, tc.want) {
				t.Fatalf("ответ %q не содержит %q", body.Response.Text, tc.want)
			}
		})
	}
}
