package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/plopezlpz/gasnow2/gas"
)

type Gas struct {
	gasServer gas.Server
	upgrader  websocket.Upgrader
}

func NewGas(infuraUrl string) *Gas {
	return &Gas{
		gasServer: gas.NewServer(infuraUrl),
		upgrader:  websocket.Upgrader{},
	}
}

func (g *Gas) GetPrice(c echo.Context) error {
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

func (g Gas) GetPriceV1(c echo.Context) error {
	gasPrices, err := g.gasServer.GetGasPrice()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Resp{
		Type: "gasprice",
		Data: gasPrices,
	})
}

type Resp struct {
	Type string  `json:"type,omitempty"`
	Data gas.Gas `json:"data,omitempty"`
}
