package presenters

import (
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/dtos"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"
	"os"
)

func RequestBodyFromParams(params entities.GenerateQRCodeParams) dtos.RequestGenerateQRCodeDTO {
	items, totalAmount := generateItems(params.Items)

	return dtos.RequestGenerateQRCodeDTO{
		Title:             "Order " + params.OrderID,
		Description:       "Order Description " + params.OrderID,
		ExternalReference: params.OrderID,
		Items:             items,
		NotificationURL:   os.Getenv("WEBHOOK_URL"), //TODO adjust this field
		TotalAmount:       FormatDecimal(totalAmount),
	}
}

func generateItems(product []entities.Item) ([]dtos.RequestGenerateQRCodeItemDTO, float64) {
	items := make([]dtos.RequestGenerateQRCodeItemDTO, len(product))
	var totalAmount float64

	for i, item := range product {
		totalAmount += item.Amount
		items[i] = dtos.RequestGenerateQRCodeItemDTO{
			Title:       item.Name,
			UnitPrice:   FormatDecimal(item.Price),
			Quantity:    item.Quantity,
			UnitMeasure: "unit",
			TotalAmount: FormatDecimal(item.Amount),
			Description: item.Description,
		}
	}
	return items, totalAmount
}
