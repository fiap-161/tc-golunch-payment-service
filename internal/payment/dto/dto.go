package dto

import (
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/entity"
	"github.com/fiap-161/tc-golunch-payment-service/internal/payment/entity/enum"
	coreentity "github.com/fiap-161/tc-golunch-payment-service/internal/shared/entity"
)

type CheckPaymentRequestDTO struct {
	Resource string `json:"resource" binding:"required"`
	Topic    string `json:"topic" binding:"required"`
}

type PaymentResponseDTO struct {
	ID      string             `json:"id"`
	OrderID string             `json:"order_id"`
	QrCode  string             `json:"qr_code"`
	Status  enum.PaymentStatus `json:"status"`
}

type PaymentListResponseDTO struct {
	Total uint                 `json:"total"`
	List  []PaymentResponseDTO `json:"list"`
}

type PaymentDAO struct {
	coreentity.Entity
	OrderID string             `json:"order_id" gorm:"not null;unique"`
	QrCode  string             `json:"qr_code" gorm:"not null"`
	Status  enum.PaymentStatus `json:"status" gorm:"not null;default:'PENDING'"`
}

func ToPaymentDAO(payment entity.Payment) PaymentDAO {
	return PaymentDAO{
		Entity:  payment.Entity,
		OrderID: payment.OrderID,
		QrCode:  payment.QrCode,
		Status:  payment.Status,
	}
}

func FromPaymentDAO(paymentDAO PaymentDAO) entity.Payment {
	return entity.Payment{
		Entity:  paymentDAO.Entity,
		OrderID: paymentDAO.OrderID,
		QrCode:  paymentDAO.QrCode,
		Status:  paymentDAO.Status,
	}
}

func EntityListFromDAOList(paymentDAOList []PaymentDAO) []entity.Payment {
	var payments []entity.Payment
	for _, dao := range paymentDAOList {
		entity := FromPaymentDAO(dao)
		payments = append(payments, entity)
	}
	return payments
}
