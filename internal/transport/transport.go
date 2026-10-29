package transport

import (
	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/alice-connect/internal/usecase"
)

type Transport struct {
	usecase         *usecase.UseCase
	protocolVersion string
}

func New(protocolVersion string) server.ServerInterface {
	return &Transport{
		usecase:         usecase.New(),
		protocolVersion: protocolVersion,
	}
}
