package transport

import (
	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/alice-connect/internal/application"
)

type Transport struct {
	app     *application.Application
	version string
}

func NewTransport(version string, app *application.Application) server.ServerInterface {
	return &Transport{
		app:     app,
		version: version,
	}
}
