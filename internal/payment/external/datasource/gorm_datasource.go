package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/dto"
	apperror "github.com/fiap-161/tc-golunch-payment-service/internal/shared/errors"
	"gorm.io/gorm"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
	Find(dest any, conds ...any) *gorm.DB
	Save(value any) *gorm.DB
}

type GormDataSource struct {
	db DB
}

func New(db DB) DataSource {
	return &GormDataSource{
		db: db,
	}
}

func (g *GormDataSource) Create(_ context.Context, payment dto.PaymentDAO) (dto.PaymentDAO, error) {
	tx := g.db.Create(&payment)
	if tx.Error != nil {
		return dto.PaymentDAO{}, tx.Error
	}

	return payment, nil
}

func (g *GormDataSource) FindByOrderID(_ context.Context, orderID string) (dto.PaymentDAO, error) {
	var payment dto.PaymentDAO

	tx := g.db.First(&payment, "order_id = ?", orderID)
	if tx.Error != nil {
		return dto.PaymentDAO{}, &apperror.NotFoundError{Msg: "Payment not found"}
	}

	return payment, nil
}

func (g *GormDataSource) Update(_ context.Context, payment dto.PaymentDAO) (dto.PaymentDAO, error) {
	tx := g.db.Save(&payment)
	if tx.Error != nil {
		return dto.PaymentDAO{}, tx.Error
	}

	return payment, nil
}

func (g *GormDataSource) GetAll(_ context.Context) ([]dto.PaymentDAO, error) {
	var payments []dto.PaymentDAO

	if err := g.db.Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
