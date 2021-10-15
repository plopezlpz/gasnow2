package infura

import (
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

func (c Client) GetLatestBlock() (Block, error) {
	var res struct {
		Result Block `json:"result"`
	}
	r, err := c.rest.R().
		SetResult(&res).
		SetBody(`{
			"id": 1337,
			"jsonrpc": "2.0",
			"method": "eth_getBlockByNumber",
			"params": ["latest", true]
		}`).
		Post(c.url)
	if err != nil {
		return Block{}, errors.Wrap(err, "getting tx from infura")
	}
	if r.IsError() {
		return Block{}, errors.Errorf("%v %v", r.Status(), string(r.Body()))
	}
	return res.Result, nil
}
