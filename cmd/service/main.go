package main

import (
	"context"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/alice-connect/internal/application"
	"github.com/Impisigmatus/alice-connect/internal/transport"
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
	const (
		version = "1.0.0"
		port    = ":8000"
		spec    = "autogen/oapi3.yaml"
		lvl     = log.LevelDebug
	)

	logger := log.New(lvl)
	server := &http.Server{
		Addr:    port,
		Handler: getRouter(logger, version, spec),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic().Err(err).Msg("Invalid service starting: ")
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
		logger.Panic().Err(err).Msg("Invalid service stopping: ")
	}
}

func getRouter(logger zerolog.Logger, version string, spec string) *chi.Mux {
	app := application.NewApplication()
	transport := transport.NewTransport(version, app)

	router := chi.NewRouter()
	router.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		http.ServeFile(w, r, spec)
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
				server.Handler(transport),
				middlewares.ContextLogger()),
			middlewares.RequestID(logger),
		),
	)

	return router
}
