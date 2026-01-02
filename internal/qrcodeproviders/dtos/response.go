package dtos

type ResponseGenerateQRCodeDTO struct {
	InStoreOrderID string `json:"in_store_order_id"`
	QRData         string `json:"qr_data"`
}

type ResponseVerifyOrderDTO struct {
	ExternalReference string `json:"external_reference"`
	OrderStatus       string `json:"order_status"`
}
