package main

import (
	"github.com/IktaS/sign/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.GET("/", handler.IndexHandler)
	app.Logger.Fatal(app.Start(":4000"))
}
