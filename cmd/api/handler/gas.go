package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/plopezlpz/gasnow2/gas"
)

type Gas struct {
	gasServer gas.Server
}

func NewGas(gasServer gas.Server) Gas {
	return Gas{
		gasServer: gasServer,
	}
}

func (g Gas) GetGasPrice(c echo.Context) error {
	gasPrices, err := g.gasServer.GetGasPrice()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, gasPrices)
}

func (g Gas) GetCurrencyPrice(c echo.Context) error {
	ethPrice, err := g.gasServer.GetCurrencyPrice("ethereum", "usd")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ethPrice)
}
