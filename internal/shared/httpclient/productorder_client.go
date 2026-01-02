package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductOrderClient struct {
	baseURL string
	client  *http.Client
}

type ProductOrder struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func NewProductOrderClient(baseURL string) *ProductOrderClient {
	return &ProductOrderClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *ProductOrderClient) FindByOrderID(ctx context.Context, orderID string) ([]ProductOrder, error) {
	url := fmt.Sprintf("%s/admin/orders/%s/products", c.baseURL, orderID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get product orders: status %d", resp.StatusCode)
	}

	var productOrders []ProductOrder
	if err := json.NewDecoder(resp.Body).Decode(&productOrders); err != nil {
		return nil, err
	}

	return productOrders, nil
}
