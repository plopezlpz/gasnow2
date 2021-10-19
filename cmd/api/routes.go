package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/plopezlpz/gasnow2/cmd/api/handler"
)

// routes sets up all the http routes of the service
func routes(e *echo.Echo, gas handler.Gas, ws handler.WS) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, json.RawMessage(`{"ok": "success"}`))
	})
	e.GET("/ws", ws.Subscribe)
	e.GET("/gasprice", gas.GetGasPrice)
	e.GET("/currencyprice", gas.GetCurrencyPrice)
}
