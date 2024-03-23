package handler

import (
	"net/http"

	"github.com/IktaS/sign/components"
	"github.com/labstack/echo/v4"
)

func FormHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Index(components.SignFile()))
}
