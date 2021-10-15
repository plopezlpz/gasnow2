package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/plopezlpz/gasnow2/cmd/api/handler"
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
	infuraProject := os.Getenv("INFURA_PROJECT")
	if infuraProject == "" {
		return fmt.Errorf("INFURA_PROJECT missing")
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000"
	}
	infuraUrl := fmt.Sprintf("%s/%s", infuraHost, infuraProject)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/ws"
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Logger.SetLevel(log.WARN)
	routes(e, handler.NewGas(infuraUrl))
	return e.Start(":" + port)
}
