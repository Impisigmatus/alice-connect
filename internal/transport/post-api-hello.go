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

// PostApiHello godoc
//
// @Router /api/hello [post]
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
func (t *Transport) PostApiHello(w http.ResponseWriter, r *http.Request) {
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

	status, err := t.usecase.SystemStatus()
	if err != nil {
		log.Error().Err(err).Msg("Не удалось получить состояние системы")
		http.Error(w, "Не удалось обработать запрос", http.StatusInternalServerError)
		return
	}

	response := &models.AliceResponse{
		Version: t.protocolVersion,
		Response: models.AliceResponseBody{
			Text:       formatSystemStatus(status),
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
	fmt.Fprintf(&text, "Привет! Go сервис на Orange Pi работает. ЦПУ: %.1f%%; ОЗУ: %.1f%% (%d из %d МБ).",
		status.CPUPercent, status.MemoryPercent, status.MemoryUsed/1024/1024, status.MemoryTotal/1024/1024)
	for _, temperature := range status.Temperatures {
		fmt.Fprintf(&text, " Температура %s: %.0f°C;", temperature.Sensor, temperature.Celsius)
	}
	return text.String()
}
