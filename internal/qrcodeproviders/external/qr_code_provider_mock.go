package external

import (
	"context"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"

	"github.com/stretchr/testify/mock"
)

type MockQRCodeProvider struct {
	mock.Mock
}

func (m *MockQRCodeProvider) GenerateQRCode(ctx context.Context, request entities.GenerateQRCodeParams) (string, error) {
	args := m.Called(ctx, request)
	return args.String(0), args.Error(1)
}

func (m *MockQRCodeProvider) CheckPayment(ctx context.Context, requestUrl string) (any, error) {
	args := m.Called(ctx, requestUrl)
	return args.Get(0), args.Error(1)
}
