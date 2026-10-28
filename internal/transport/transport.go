package transport

import (
	"io"
	"net/http"

	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/service_core/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type Transport struct{}

func NewTransport() server.ServerInterface {
	return &Transport{}
}

// Set godoc
//
// @Router /api/hello [post]
// @Summary Ответ на сообщение Алисы
// @Description Возвращает приветствие и идентификаторы полученной сессии.
// @Accept json
// @Produce json
//
// @Tags APIs
//
// @Param request body AliceRequest true "Запрос Алисы"
// @Success 200 {object} AliceResponse "Ответ Алисе"
// @Failure 400 {string} string "Некорректный JSON-запрос или запрос Алисы"
// @Failure 500 {string} string "Не удалось обработать запрос"
func (transport *Transport) PostApiHello(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		http.Error(w, "Невалидный логгер", http.StatusInternalServerError)
		return
	}

	var req AliceRequest
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

	response := AliceResponse{
		Version: "1.0",
		Response: AliceResponseBody{
			Text:       "Привет! Go сервис на Orange Pi работает.",
			EndSession: false,
		},
		Session: AliceResponseSession{
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
