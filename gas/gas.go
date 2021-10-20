package gas

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/plopezlpz/gasnow2/coingecko"
	"github.com/plopezlpz/gasnow2/infura"
	"github.com/plopezlpz/gasnow2/numb"
)

type Server struct {
	infuraClient infura.Client
	geckoClient  coingecko.Client
	cache        *cache.Cache
}

func NewServer(infuraClient infura.Client, geckoClient coingecko.Client) Server {
	return Server{
		infuraClient: infuraClient,
		geckoClient:  geckoClient,
		cache:        cache.New(5*time.Second, -1),
	}
}

func (s Server) GetGasPrice() (Gas, error) {
	if item, exp, found := s.cache.GetWithExpiration("gasprice"); found && time.Now().Before(exp) {
		gasCached, ok := item.(Gas)
		if !ok {
			return Gas{}, errors.New("wrong gasprice type in cache")
		}
		return gasCached, nil
	}
	block, err := s.infuraClient.GetPendingBlock()
	if err != nil {
		return Gas{}, err
	}
	if len(block.Transactions) < 1 {
		return Gas{}, fmt.Errorf("no transactions in block %v", numb.ToDecimal(block.Number, 0))
	}
	sort.Slice(block.Transactions, func(i, j int) bool {
		return numb.ToDecimal(block.Transactions[j].Gasprice, 0).LessThan(numb.ToDecimal(block.Transactions[i].Gasprice, 0))
	})

	timestamp, err := numb.ToInt64(block.Timestamp)
	if err != nil {
		return Gas{}, err
	}

	rapid, err := numb.ToGwei(block.Transactions[len(block.Transactions)/2].Gasprice)
	if err != nil {
		return Gas{}, err
	}
	fast, err := numb.ToGwei(block.Transactions[len(block.Transactions)-1].Gasprice)
	if err != nil {
		return Gas{}, err
	}
	newGas := Gas{
		GasPrices: GasPrice{
			Rapid: rapid,
			Fast:  fast,
		},
		Timestamp: timestamp,
	}
	s.cache.Set("gasprice", newGas, 0)
	return newGas, nil
}

func (s Server) GetCurrencyPrice(currency, vsCurrency string) (coingecko.Price, error) {
	if item, exp, found := s.cache.GetWithExpiration("currencyprice"); found && time.Now().Before(exp) {
		currencyPriceCached, ok := item.(coingecko.Price)
		if !ok {
			return coingecko.Price{}, errors.New("wrong currencyprice type in cache")
		}
		return currencyPriceCached, nil
	}
	price, err := s.geckoClient.GetCurrencyPrice(currency, vsCurrency)
	if err != nil {
		return coingecko.Price{}, err
	}
	s.cache.Set("currencyprice", price, 30*time.Second)
	return price, nil
}
