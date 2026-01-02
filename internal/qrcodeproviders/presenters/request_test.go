package presenters

import (
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/dtos"
	"github.com/fiap-161/tc-golunch-payment-service/internal/qrcodeproviders/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBodyFromParams(t *testing.T) {
	type args struct {
		params entities.GenerateQRCodeParams
	}
	tests := []struct {
		name string
		args args
		want dtos.RequestGenerateQRCodeDTO
	}{
		{
			name: "Given a valid params with multiple items",
			args: args{
				params: entities.GenerateQRCodeParams{
					OrderID: "orderId",
					Items: []entities.Item{
						{
							ID:          "itemId",
							Name:        "itemName",
							Price:       12.3,
							Description: "itemDescription",
							Quantity:    1,
							Amount:      12.3,
						},
						{
							ID:          "itemId2",
							Name:        "itemName2",
							Price:       45.6,
							Description: "itemDescription2",
							Quantity:    2,
							Amount:      91.2,
						},
					},
				},
			},
			want: dtos.RequestGenerateQRCodeDTO{
				Title:             "Order orderId",
				Description:       "Order Description orderId",
				ExternalReference: "orderId",
				NotificationURL:   "",
				TotalAmount:       103.5,
				Items: []dtos.RequestGenerateQRCodeItemDTO{
					{Title: "itemName",
						UnitPrice:   12.3,
						Quantity:    1,
						UnitMeasure: "unit",
						TotalAmount: 12.3,
						SkuNumber:   "",
						Category:    "",
						Description: "itemDescription",
					},
					{Title: "itemName2",
						UnitPrice:   45.6,
						Quantity:    2,
						UnitMeasure: "unit",
						TotalAmount: 91.2,
						SkuNumber:   "",
						Category:    "",
						Description: "itemDescription2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequestBodyFromParams(tt.args.params)
			assert.Equal(t, tt.want, got)
		})
	}
}
