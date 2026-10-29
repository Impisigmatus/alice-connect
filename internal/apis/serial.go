package apis

import (
	"github.com/Impisigmatus/alice-connect/internal/config"
	"github.com/go-resty/resty/v2"
)

type SerialDriverAPI struct {
	client *resty.Client
}

func NewSerialDriverAPI(cfg config.SerialDriverAPIConfig) *SerialDriverAPI {
	return &SerialDriverAPI{
		client: resty.New().
			SetBaseURL(cfg.Hostname).
			SetTimeout(cfg.Timeout).
			SetHeader("Accept", cfg.ContentType).
			SetRetryCount(cfg.RetryCount),
	}
}

func (api *SerialDriverAPI) GetSerialLatest() (string, error) {
	const endpoint = "/api/serial/latest"
	resp, err := api.client.R().Get(endpoint)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", &resty.ResponseError{Response: resp}
	}

	return resp.String(), nil
}
