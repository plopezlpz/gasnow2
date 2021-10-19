package coingecko

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type Client struct {
	url  string
	rest *resty.Client
}

func NewClient(url string) Client {
	return Client{
		url:  url,
		rest: resty.New(),
	}
}

func (c Client) GetCurrencyPrice(currency, vsCurrency string) (Price, error) {
	res := make(map[string]Price, 4)
	r, err := c.rest.R().
		SetResult(&res).
		SetHeader("accept", "application/json").
		SetQueryParam("ids", currency).
		SetQueryParam("vs_currencies", vsCurrency).
		SetQueryParam("include_last_updated_at", "true").
		Get(fmt.Sprintf("%s/simple/price", c.url))
	if err != nil {
		return Price{}, errors.Wrap(err, "getting price from coingecko")
	}
	if r.IsError() {
		return Price{}, errors.Errorf("%v %v", r.Status(), string(r.Body()))
	}
	data, ok := res[currency]
	if !ok {
		return Price{}, errors.Errorf("%v not returned", currency)
	}
	data.Currency = currency
	return data, nil
}
