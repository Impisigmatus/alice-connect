package transport

import (
	"net/http"

	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/service_core/log"
	"github.com/rs/zerolog"
)

type Transport struct{}

func NewTransport() server.ServerInterface {
	return &Transport{}
}

// Set godoc
//
// @Router /api/hello [get]
// @Summary Получение приветствия
// @Description Возвращает фиксированное текстовое сообщение «Привет от Golang!».
// @Produce plain
//
// @Tags APIs
//
// @Success 200 {string} string "Привет от Golang!"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiHello(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		http.Error(w, "Невалидный логгер", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("Привет от Golang!")); err != nil {
		log.Error().Err(err).Msg("Не удалось записать ответ")
	}
}
