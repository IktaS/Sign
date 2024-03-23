package handler

import (
	"net/http"

	"github.com/IktaS/sign/components"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func IndexHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Index(components.Form()))
}
