package gas

import (
	"errors"
	"sort"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/plopezlpz/gasnow2/infura"
	"github.com/plopezlpz/gasnow2/numb"
)

type Server struct {
	infuraClient infura.Client
	cache        *cache.Cache
}

func NewServer(infuraUrl string) Server {
	return Server{
		infuraClient: infura.NewClient(infuraUrl),
		cache:        cache.New(7*time.Second, -1),
	}
}

func (s Server) GetGasPrice() (Gas, error) {
	if item, exp, found := s.cache.GetWithExpiration("gasprice"); found && time.Now().Before(exp) {
		gasCached, ok := item.(Gas)
		if !ok {
			return Gas{}, errors.New("wrong type in cache")
		}
		return gasCached, nil
	}
	block, err := s.infuraClient.GetLatestBlock()
	if err != nil {
		return Gas{}, err
	}
	sort.Slice(block.Transactions, func(i, j int) bool {
		return numb.ToDecimal(block.Transactions[j].Gasprice, 0).LessThan(numb.ToDecimal(block.Transactions[i].Gasprice, 0))
	})

	timestamp, err := numb.ToTimestamp(block.Timestamp)
	if err != nil {
		return Gas{}, err
	}
	newGas := Gas{
		GasPrices: GasPrice{
			Rapid: numb.ToDecimal(block.Transactions[len(block.Transactions)/2].Gasprice, -9),
			Fast:  numb.ToDecimal(block.Transactions[len(block.Transactions)-1].Gasprice, -9),
		},
		Timestamp: timestamp,
	}
	s.cache.Set("gasprice", newGas, 0)
	return newGas, nil
}
