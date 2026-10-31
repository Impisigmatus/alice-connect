package transport

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Impisigmatus/alice-connect/internal/models"
	"github.com/Impisigmatus/service_core/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

// PostApiAliceWebhook godoc
//
// @Router /api/alice/webhook [post]
// @Summary Ответ на сообщение Алисы
// @Description Возвращает приветствие, загрузку ЦПУ и ОЗУ, доступную температуру и идентификаторы сессии.
// @Accept json
// @Produce json
//
// @Tags APIs
//
// @Param request body models.AliceRequest true "Запрос Алисы"
// @Success 200 {object} models.AliceResponse "Ответ Алисе"
// @Failure 400 {string} string "Некорректный JSON-запрос или запрос Алисы"
// @Failure 500 {string} string "Не удалось обработать запрос"
func (t *Transport) PostApiAliceWebhook(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		http.Error(w, "Невалидный логгер", http.StatusInternalServerError)
		return
	}

	var req models.AliceRequest
	requestBody, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		http.Error(w, "Некорректный JSON-запрос", http.StatusBadRequest)
		return
	}

	if err := jsoniter.Unmarshal(requestBody, &req); err != nil {
		http.Error(w, "Некорректный JSON-запрос", http.StatusBadRequest)
		return
	}

	if req.Session.SessionID == "" {
		http.Error(w, "Некорректный запрос Алисы", http.StatusBadRequest)
		return
	}

	var result string
	command := strings.Join(strings.Fields(strings.ToLower(req.Request.Command)), " ")
	switch command {
	case "какая загрузка процессора", "загрузка процессора":
		status, err := t.usecase.SystemStatus()
		if err != nil {
			log.Error().Err(err).Msg("Не удалось получить состояние системы")
			http.Error(w, "Не удалось обработать запрос", http.StatusInternalServerError)
			return
		}
		result = formatSystemStatus(status)
	case "сожми манипулятор", "сжать манипулятор":
		data, err := t.usecase.SerialCommand("TAKE")
		if err != nil {
			log.Error().Err(err).Msg("Не удалось сжать манипулятор")
			http.Error(w, "Не удалось выполнить команду манипулятора", http.StatusInternalServerError)
			return
		}
		result = data
	case "разожми манипулятор", "разжать манипулятор":
		data, err := t.usecase.SerialCommand("RELEASE")
		if err != nil {
			log.Error().Err(err).Msg("Не удалось разжать манипулятор")
			http.Error(w, "Не удалось выполнить команду манипулятора", http.StatusInternalServerError)
			return
		}
		result = data
	default:
		result = "Привет! Go сервис на Orange Pi работает."
	}

	response := &models.AliceResponse{
		Version: t.protocolVersion,
		Response: models.AliceResponseBody{
			Text:       result,
			EndSession: true,
		},
		Session: models.AliceResponseSession{
			SessionID: req.Session.SessionID,
			MessageID: req.Session.MessageID,
		},
	}
	responseBody, err := jsoniter.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Не удалось сформировать ответ")
		http.Error(w, "Не удалось сформировать ответ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(responseBody); err != nil {
		log.Error().Err(err).Msg("Не удалось записать ответ")
	}
}

func formatSystemStatus(status *models.SystemStatus) string {
	var text strings.Builder
	fmt.Fprintf(&text, "ЦПУ: %.1f%%; ОЗУ: %.1f%% (%d из %d МБ).",
		status.CPUPercent, status.MemoryPercent, status.MemoryUsed/1024/1024, status.MemoryTotal/1024/1024)
	for _, temperature := range status.Temperatures {
		fmt.Fprintf(&text, " Температура %s: %.0f°C;", temperature.Sensor, temperature.Celsius)
	}
	return text.String()
}
