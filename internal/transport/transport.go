package transport

import (
	"github.com/Impisigmatus/alice-connect/autogen/server"
	"github.com/Impisigmatus/alice-connect/internal/usecase"
)

type Transport struct {
	usecase         *usecase.UseCase
	protocolVersion string
}

func New(protocolVersion string, usecase *usecase.UseCase) server.ServerInterface {
	return &Transport{
		usecase:         usecase,
		protocolVersion: protocolVersion,
	}
}
