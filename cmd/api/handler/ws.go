package handler

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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

func (g WS) Subscribe(c echo.Context) error {
	ws, err := g.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		gasPrices, err := g.gasServer.GetGasPrice()
		if err != nil {
			c.Logger().Error(err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Write
		if err := ws.WriteJSON(Resp{
			Type: "gasprice",
			Data: gasPrices,
		}); err != nil {
			c.Logger().Warnf("connection closed: %v", err)
			return nil
		}
		time.Sleep(5 * time.Second)
	}
}

type Resp struct {
	Type string  `json:"type,omitempty"`
	Data gas.Gas `json:"data,omitempty"`
}
