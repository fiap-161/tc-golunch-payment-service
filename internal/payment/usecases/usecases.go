package usecases

import (
	"context"
	"fmt"

	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/external"

	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/entity"
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/entity/enum"
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/gateway"
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/interfaces"
	apperror "github.com/fiap-161/tc-golunch-payment-service/internal/shared/errors"
)

type UseCases struct {
	paymentGateway      *gateway.Gateway
	qrCodeProvider      external.QRCodeProvider
	productService      interfaces.ProductService
	productOrderService interfaces.ProductOrderService
	orderService        interfaces.OrderService
}

func Build(
	paymentGateway *gateway.Gateway,
	qrCodeProvider external.QRCodeProvider,
	productService interfaces.ProductService,
	productOrderService interfaces.ProductOrderService,
	orderService interfaces.OrderService,
) *UseCases {
	return &UseCases{
		paymentGateway:      paymentGateway,
		qrCodeProvider:      qrCodeProvider,
		productService:      productService,
		productOrderService: productOrderService,
		orderService:        orderService,
	}
}

func (u *UseCases) CreateByOrderID(ctx context.Context, orderID string) (entity.Payment, error) {
	productOrders, productOrderErr := u.productOrderService.FindByOrderID(ctx, orderID)
	if productOrderErr != nil {
		return entity.Payment{}, productOrderErr
	}

	var productIDs []string
	for _, po := range productOrders {
		productIDs = append(productIDs, po.ProductID)
	}

	products, productsErr := u.productService.FindByIDs(ctx, productIDs)
	if productsErr != nil {
		return entity.Payment{}, productsErr
	}

	var items []entities.Item
	for _, po := range productOrders {
		for _, product := range products {
			if po.ProductID == product.ID {
				items = append(items, entities.Item{
					ID:          product.ID,
					Name:        product.Name,
					Price:       product.Price,
					Description: product.Name, // Usar Name como Description já que não tem Description
					Quantity:    po.Quantity,
					Amount:      product.Price * float64(po.Quantity),
				})
				break
			}
		}
	}

	qrCode, qrCodeErr := u.qrCodeProvider.GenerateQRCode(ctx, entities.GenerateQRCodeParams{
		OrderID: orderID,
		Items:   items,
	})
	if qrCodeErr != nil {
		return entity.Payment{}, qrCodeErr
	}

	var payment entity.Payment
	createdPayment, createErr := u.paymentGateway.Create(ctx, payment.Build(orderID, qrCode))
	if createErr != nil {
		return entity.Payment{}, createErr
	}

	return createdPayment, nil
}

func (u *UseCases) CheckPayment(ctx context.Context, requestUrl string) (interface{}, error) {
	if requestUrl == "" {
		return nil, &apperror.ValidationError{Msg: "Request URL is required"}
	}

	response, err := u.qrCodeProvider.CheckPayment(ctx, requestUrl)
	if err != nil {
		return nil, fmt.Errorf("error checking payment: %w", err)
	}

	payment, paymentErr := u.paymentGateway.FindByOrderID(ctx, response.ExternalReference)
	if paymentErr != nil {
		return nil, paymentErr
	}
	if response.OrderStatus == "paid" {
		payment.Status = enum.PaymentStatusApproved
		_, updateErr := u.paymentGateway.Update(ctx, payment)
		if updateErr != nil {
			return nil, updateErr
		}

		order, orderErr := u.orderService.FindByID(ctx, response.ExternalReference)
		if orderErr != nil {
			return nil, orderErr
		}

		order.Status = "RECEIVED"
		_, updateOrderErr := u.orderService.Update(ctx, order)
		if updateOrderErr != nil {
			return nil, updateOrderErr
		}
	}

	return response, nil
}
