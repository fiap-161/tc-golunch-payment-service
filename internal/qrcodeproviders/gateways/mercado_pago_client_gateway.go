package gateways

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MercadoPagoClientRest struct {
	client *resty.Client
}

func (r *MercadoPagoClientRest) Get(url string, result interface{}) (*resty.Response, error) {
	resp, err := r.client.R().
		SetResult(result).
		Get(url)

	if err != nil {
		return resp, err
	}

	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return resp, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp, nil
}

func (r *MercadoPagoClientRest) Post(url string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := r.client.R().
		SetBody(body).
		Post(url)

	if err != nil {
		return resp, err
	}

	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return resp, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp, nil
}
