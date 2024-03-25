package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/IktaS/sign/components"
	"github.com/IktaS/sign/service"
	"github.com/labstack/echo/v4"
)

func VerifyHandler(s *service.SignService) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		if id == "" {
			c.NoContent(http.StatusNotFound)
			return nil
		}
		info, err := s.GetSignatureInfo(c.Request().Context(), id)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		return Render(c, http.StatusOK, components.VerifyFile(id, info.Filename, info.Fullname, info.CreatedAt.Format(time.RFC3339)))
	}
}

func VerifyFileHandler(s *service.SignService) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.FormValue("id")
		if id == "" {
			c.NoContent(http.StatusNotFound)
			return nil
		}
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		b, err := io.ReadAll(src)
		if err != nil {
			return err
		}

		isVerified, err := s.VerifyFileHash(c.Request().Context(), id, b)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		if isVerified {
			c.NoContent(http.StatusOK)
		} else {
			c.NoContent(http.StatusUnauthorized)
		}
		return nil
	}
}
