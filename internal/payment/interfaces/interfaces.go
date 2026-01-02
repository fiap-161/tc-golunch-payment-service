package interfaces

import (
	"context"

	"github.com/fiap-161/tc-golunch-payment-service/internal/shared/httpclient"
)

type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]httpclient.Product, error)
}

type ProductOrderService interface {
	FindByOrderID(ctx context.Context, orderID string) ([]httpclient.ProductOrder, error)
}

type OrderService interface {
	FindByID(ctx context.Context, orderID string) (httpclient.Order, error)
	Update(ctx context.Context, order httpclient.Order) (httpclient.Order, error)
}
