package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/dto"
)

type DataSource interface {
	Create(ctx context.Context, payment dto.PaymentDAO) (dto.PaymentDAO, error)
	FindByOrderID(ctx context.Context, orderID string) (dto.PaymentDAO, error)
	Update(ctx context.Context, payment dto.PaymentDAO) (dto.PaymentDAO, error)
	GetAll(ctx context.Context) ([]dto.PaymentDAO, error)
}
