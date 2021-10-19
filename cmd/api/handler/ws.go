package handler

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/plopezlpz/gasnow2/coingecko"
	"github.com/plopezlpz/gasnow2/gas"
)

type WS struct {
	upgrader  *websocket.Upgrader
	gasServer gas.Server
}

func NewWS(upgrader *websocket.Upgrader, gasServer gas.Server) WS {
	return WS{
		upgrader:  upgrader,
		gasServer: gasServer,
	}
}

func (w WS) Subscribe(c echo.Context) error {
	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(c.Request().Context())
	wg := sync.WaitGroup{}
	wg.Add(3)

	gasPrice := w.subscribeToGasPrice(ctx, c.Logger(), &wg)
	currencyPrice := w.subscribeToCurrencyPrice(ctx, c.Logger(), &wg)

	go func() {
		defer wg.Done()
		for {
			select {
			case gp := <-gasPrice:
				if err := ws.WriteJSON(gp); err != nil {
					c.Logger().Error(err)
					cancel()
					return
				}
			case cp := <-currencyPrice:
				if err := ws.WriteJSON(cp); err != nil {
					c.Logger().Error(err)
					cancel()
					return
				}
			}
		}
	}()
	wg.Wait()
	return nil
}

func (w WS) subscribeToGasPrice(ctx context.Context, logger echo.Logger, wg *sync.WaitGroup) <-chan RespGas {
	out := make(chan RespGas)
	var lastUpdate int64 = 0

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default: // avoid blocking
			}
			gasPrices, err := w.gasServer.GetGasPrice()
			if err != nil {
				logger.Error(err)
				time.Sleep(5 * time.Second)
				continue
			}
			if gasPrices.Timestamp == lastUpdate {
				time.Sleep(5 * time.Second)
				continue
			}
			lastUpdate = gasPrices.Timestamp
			out <- RespGas{
				Type: "gasprice",
				Data: gasPrices,
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return out

}

func (w WS) subscribeToCurrencyPrice(ctx context.Context, logger echo.Logger, wg *sync.WaitGroup) <-chan RespCurrency {
	out := make(chan RespCurrency)
	var lastUpdate int64 = 0

	go func() {
		defer wg.Done()
		for {
			// Check if any error occurred in any other gorouties:
			select {
			case <-ctx.Done():
				return
			default: // avoid blocking
			}
			currencyPrice, err := w.gasServer.GetCurrencyPrice("ethereum", "usd")
			if err != nil {
				logger.Error(err)
				close(out)
				return
			}
			if currencyPrice.LastUpdatedAt == lastUpdate {
				time.Sleep(5 * time.Second)
				continue
			}
			lastUpdate = currencyPrice.LastUpdatedAt
			out <- RespCurrency{
				Type: "currencyprice",
				Data: currencyPrice,
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return out
}

type RespGas struct {
	Type string  `json:"type,omitempty"`
	Data gas.Gas `json:"data,omitempty"`
}

type RespCurrency struct {
	Type string          `json:"type,omitempty"`
	Data coingecko.Price `json:"data,omitempty"`
}
