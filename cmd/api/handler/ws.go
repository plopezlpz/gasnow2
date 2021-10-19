package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/plopezlpz/gasnow2/coingecko"
	"github.com/plopezlpz/gasnow2/gas"
	"golang.org/x/sync/errgroup"
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

	g, ctx := errgroup.WithContext(c.Request().Context())
	g.Go(func() error {
		return w.subscribeToGasPrice(ctx, c.Logger(), ws)
	})
	g.Go(func() error {
		return w.subscribeToCurrencyPrice(ctx, c.Logger(), ws)
	})
	if err := g.Wait(); err != nil {
		c.Logger().Error(err)
	}
	return nil
}

func (w WS) subscribeToGasPrice(ctx context.Context, logger echo.Logger, ws *websocket.Conn) error {
	for {
		// Check if any error occurred in any other gorouties:
		select {
		case <-ctx.Done():
			return fmt.Errorf("stopping %v", "gaspriceSubscription")
		default: // to avoid blocking
		}
		gasPrices, err := w.gasServer.GetGasPrice()
		if err != nil {
			logger.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Write
		if err := ws.WriteJSON(RespGas{
			Type: "gasprice",
			Data: gasPrices,
		}); err != nil {
			return errors.Wrap(err, "connection closed")
		}
		time.Sleep(5 * time.Second)
	}
}

func (w WS) subscribeToCurrencyPrice(ctx context.Context, logger echo.Logger, ws *websocket.Conn) error {
	var lastUpdate int64 = 0
	for {
		// Check if any error occurred in any other gorouties:
		select {
		case <-ctx.Done():
			return fmt.Errorf("stopping %v", "currencyPriceSubscription")
		default: // to avoid blocking
		}
		currencyPrice, err := w.gasServer.GetCurrencyPrice("ethereum", "usd")
		if err != nil {
			return errors.Wrap(err, "getting currency price")
		}
		if currencyPrice.LastUpdatedAt == lastUpdate {
			time.Sleep(5 * time.Second)
			continue
		}
		lastUpdate = currencyPrice.LastUpdatedAt
		// Write
		if err := ws.WriteJSON(RespCurrency{
			Type: "currencyprice",
			Data: currencyPrice,
		}); err != nil {
			return errors.Wrap(err, "connection closed")
		}
		time.Sleep(5 * time.Second)
	}
}

type RespGas struct {
	Type string  `json:"type,omitempty"`
	Data gas.Gas `json:"data,omitempty"`
}

type RespCurrency struct {
	Type string          `json:"type,omitempty"`
	Data coingecko.Price `json:"data,omitempty"`
}
