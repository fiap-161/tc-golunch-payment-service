package external

import "github.com/go-resty/resty/v2"

type MercadoPagoClient interface {
	Get(url string, result interface{}) (*resty.Response, error)
	Post(url string, body interface{}, result interface{}) (*resty.Response, error)
}
