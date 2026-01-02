package external

import (
	"context"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/dtos"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request entities.GenerateQRCodeParams) (string, error)
	CheckPayment(ctx context.Context, requestUrl string) (dtos.ResponseVerifyOrderDTO, error)
}
