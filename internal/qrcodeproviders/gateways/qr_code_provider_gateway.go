package gateways

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/dtos"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"
	external2 "github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/external"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/presenters"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"github.com/fiap-161/tc-golunch-payment-service/internal/shared"
)

func GetSellerUserID() string {
	return os.Getenv("MERCADO_PAGO_SELLER_APP_USER_ID")
}

func GetExternalPosID() string {
	return os.Getenv("MERCADO_PAGO_EXTERNAL_POS_ID")
}

type MercadoPagoClient struct {
	client external2.MercadoPagoClient
}

func New() external2.QRCodeProvider {
	return &MercadoPagoClient{
		client: getClient(),
	}
}

func getClient() external2.MercadoPagoClient {
	return &MercadoPagoClientRest{
		client: resty.New().
			SetBaseURL(viper.GetString(shared.MercadoPagoHost)).
			SetHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"),
			}),
	}
}

func (m *MercadoPagoClient) GenerateQRCode(_ context.Context, params entities.GenerateQRCodeParams) (string, error) {
	requestBody := presenters.RequestBodyFromParams(params)

	pathParams := []shared.BuildPathParam{
		{
			Key:   "user_id",
			Value: GetSellerUserID(),
		},
		{
			Key:   "external_pos_id",
			Value: GetExternalPosID(),
		},
	}
	resolvedPath, err := shared.BuildPath(viper.GetString(shared.MercadoPagoQRCodePath), pathParams)
	if err != nil {
		return "", err
	}

	var responseDTO dtos.ResponseGenerateQRCodeDTO
	res, reqErr := m.client.Post(resolvedPath, requestBody, &responseDTO)

	if res != nil && res.IsError() {
		fmt.Println(" ")
		fmt.Println("❌ Erro na chamada MercadoPago")
		fmt.Println("Status Code:", res.StatusCode())
		fmt.Println("Request URL:", res.Request.URL)
		fmt.Println("Response Body:", string(res.Body()))
		fmt.Println("*** Notification URL: ", requestBody.NotificationURL)
		fmt.Println(" ")

		return "", errors.New("error in request, endpoint called: " + res.Request.URL)
	}

	if reqErr != nil {
		return "", reqErr
	}

	return responseDTO.QRData, nil
}

func (m *MercadoPagoClient) CheckPayment(_ context.Context, requestUrl string) (dtos.ResponseVerifyOrderDTO, error) {
	fmt.Println(requestUrl)

	var responseDTO dtos.ResponseVerifyOrderDTO
	res, reqErr := m.client.Get(requestUrl, &responseDTO)

	if res != nil && res.IsError() {
		return dtos.ResponseVerifyOrderDTO{}, errors.New("error in request, endpoint called: " + res.Request.URL + "\n")
	}
	if reqErr != nil {
		return dtos.ResponseVerifyOrderDTO{}, reqErr
	}

	return responseDTO, nil
}
