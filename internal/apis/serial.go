package apis

import (
	"github.com/Impisigmatus/alice-connect/internal/config"
	"github.com/Impisigmatus/alice-connect/internal/models"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
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

func (api *SerialDriverAPI) PostSerialCommand(cmd string) (string, error) {
	const endpoint = "/api/serial/command"
	resp, err := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(models.SerialDriverRequest{Data: cmd}).
		Post(endpoint)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", &resty.ResponseError{Response: resp}
	}

	var result models.SerialDriverResponse
	if err := jsoniter.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	return result.Data, nil
}
