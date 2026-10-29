package main

import (
	"context"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/alice-connect/internal/config"
	"github.com/Impisigmatus/alice-connect/internal/transport"
	"github.com/Impisigmatus/alice-connect/internal/usecase"
	"github.com/Impisigmatus/service_core/log"
	"github.com/Impisigmatus/service_core/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Alice Connect API
// @version 1.0.0
// @description Ответы на сообщения Алисы
func main() {
	logger := log.New(log.LevelPanic)
	cfg, err := config.Load()
	if err != nil {
		logger.Panic().Err(err).Msg("Invalid config")
	}

	logger = log.New(cfg.LogLevel)
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: getRouter(logger, cfg),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic().Err(err).Msg("Invalid service starting")
		}
		logger.Info().Msg("Service stopped")
	}()
	logger.Info().Msg("Service started")

	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGABRT,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-channel

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Panic().Err(err).Msg("Invalid service stopping")
	}
}

func getRouter(logger zerolog.Logger, cfg config.Config) *chi.Mux {
	usecase := usecase.New(cfg.SerialDriverAPI)
	srv := transport.New(cfg.AliceProtocolVersion, usecase)

	router := chi.NewRouter()
	router.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		http.ServeFile(w, r, cfg.SpecPath)
	})
	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/openapi.yaml")))

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/api/*",
		middlewares.Use(
			middlewares.Use(
				server.Handler(srv),
				middlewares.ContextLogger()),
			middlewares.RequestID(logger),
		),
	)

	return router
}
