package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/plopezlpz/gasnow2/cmd/api/handler"
	"github.com/plopezlpz/gasnow2/coingecko"
	"github.com/plopezlpz/gasnow2/gas"
	"github.com/plopezlpz/gasnow2/infura"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR main: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		errors.Wrap(err, "loading env")
	}
	infuraHost := os.Getenv("INFURA_URL")
	if infuraHost == "" {
		return fmt.Errorf("INFURA_URL missing")
	}
	geckoUrl := os.Getenv("GECKO_URL")
	if geckoUrl == "" {
		return fmt.Errorf("GECKO_URL missing")
	}
	infuraProject := os.Getenv("INFURA_PROJECT")
	if infuraProject == "" {
		return fmt.Errorf("INFURA_PROJECT missing")
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000"
	}
	logLevel := parseLogLvl(os.Getenv("LOG_LEVEL"))
	infuraUrl := fmt.Sprintf("%s/%s", infuraHost, infuraProject)

	e := echo.New()
	e.Use(middleware.CORS()) // TODO restrict
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/ws"
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Logger.SetLevel(logLevel)

	wsUpgrader := websocket.Upgrader{}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true // TODO restrict
	}

	gasServer := gas.NewServer(infura.NewClient(infuraUrl), coingecko.NewClient(geckoUrl))
	routes(e, handler.NewGas(gasServer), handler.NewWS(&wsUpgrader, gasServer))
	return e.Start(":" + port)
}

func parseLogLvl(val string) log.Lvl {
	val = strings.ToUpper(val)
	switch val {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OFF":
		return log.OFF
	default:
		return log.ERROR
	}
}
