package dtos

type RequestGenerateQRCodeDTO struct {
	Title             string                         `json:"title"`
	Description       string                         `json:"description"`
	ExternalReference string                         `json:"external_reference"`
	Items             []RequestGenerateQRCodeItemDTO `json:"items"`
	NotificationURL   string                         `json:"notification_url"`
	TotalAmount       float64                        `json:"total_amount"`
}

type RequestGenerateQRCodeItemDTO struct {
	Title       string  `json:"title"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	UnitMeasure string  `json:"unit_measure"`
	TotalAmount float64 `json:"total_amount"`
	SkuNumber   string  `json:"sku_number,omitempty"`
	Category    string  `json:"category,omitempty"`
	Description string  `json:"description,omitempty"`
}

func (r RequestGenerateQRCodeDTO) GetItems() []RequestGenerateQRCodeItemDTO {
	return r.Items
}
