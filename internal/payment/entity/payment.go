package entity

import (
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/entity/enum"
	"github.com/fiap-161/tc-golunch-payment-service/internal/shared/entity"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	entity.Entity
	OrderID string             `json:"order_id" gorm:"not null;unique"`
	QrCode  string             `json:"qr_code" gorm:"not null"`
	Status  enum.PaymentStatus `json:"status" gorm:"not null;default:'PENDING'"`
}

func (p Payment) Build(orderID, qrCode string) Payment {
	now := time.Now()

	return Payment{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderID: orderID,
		QrCode:  qrCode,
		Status:  enum.PaymentStatusPending,
	}
}
