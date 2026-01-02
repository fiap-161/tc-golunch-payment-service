package entities

type GenerateQRCodeParams struct {
	OrderID string
	Items   []Item
}

type Item struct {
	ID          string
	Name        string
	Price       float64
	Description string
	Quantity    int
	Amount      float64
}
