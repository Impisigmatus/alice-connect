package usecase

import (
	"github.com/Impisigmatus/alice-connect/internal/apis"
	"github.com/Impisigmatus/alice-connect/internal/config"
)

type UseCase struct {
	serialDriverAPI *apis.SerialDriverAPI
}

func New(cfg config.SerialDriverAPIConfig) *UseCase {
	return &UseCase{
		serialDriverAPI: apis.NewSerialDriverAPI(cfg),
	}
}
