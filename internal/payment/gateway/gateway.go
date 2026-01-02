package gateway

import (
	"context"
	"errors"

	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/dto"
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/entity"
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/external/datasource"
	apperror "github.com/fiap-161/tc-golunch-payment-service/internal/shared/errors"
)

type Gateway struct {
	datasource datasource.DataSource
}

func Build(datasource datasource.DataSource) *Gateway {
	return &Gateway{
		datasource: datasource,
	}
}

func (g *Gateway) Create(c context.Context, payment entity.Payment) (entity.Payment, error) {
	var paymentDAO = dto.ToPaymentDAO(payment)
	created, err := g.datasource.Create(c, paymentDAO)

	if err != nil {
		return entity.Payment{}, &apperror.InternalError{Msg: err.Error()}
	}

	return dto.FromPaymentDAO(created), nil
}

func (g *Gateway) FindByOrderID(c context.Context, orderID string) (entity.Payment, error) {
	found, err := g.datasource.FindByOrderID(c, orderID)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return entity.Payment{}, notFoundErr
		}
		return entity.Payment{}, &apperror.InternalError{Msg: "Unexpected error"}
	}

	return dto.FromPaymentDAO(found), nil
}

func (g *Gateway) Update(c context.Context, payment entity.Payment) (entity.Payment, error) {
	paymentDAO := dto.ToPaymentDAO(payment)
	updated, err := g.datasource.Update(c, paymentDAO)

	if err != nil {
		return entity.Payment{}, &apperror.InternalError{Msg: err.Error()}
	}

	return dto.FromPaymentDAO(updated), nil
}
