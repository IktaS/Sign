package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IktaS/sign/service"
	"github.com/labstack/echo/v4"
)

func SignHandler(s *service.SignService) func(c echo.Context) error {
	return func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		locationWidth := c.FormValue("qr-location-width")
		locationX, err := strconv.Atoi(locationWidth)
		if err != nil {
			return err
		}
		locationHeight := c.FormValue("qr-location-height")
		locationY, err := strconv.Atoi(locationHeight)
		if err != nil {
			return err
		}
		qrSizeStr := c.FormValue("qr-size")
		qrSize, err := strconv.Atoi(qrSizeStr)
		if err != nil {
			return err
		}
		var qrPage *int
		qrPageStr := c.FormValue("qr-page")
		if len(qrPageStr) > 0 {
			p, err := strconv.Atoi(qrPageStr)
			if err != nil {
				return err
			}
			qrPage = &p
		}

		var isAllPage *bool
		allPageStr := c.FormValue("all-page")
		if len(allPageStr) > 0 {
			t := true
			isAllPage = &t
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

		req := service.SignRequest{
			Username:            username,
			Password:            password,
			LocationXPercentage: locationX,
			LocationYPercentage: locationY,
			QRSize:              qrSize,
			QRPage:              qrPage,
			IsAllPage:           isAllPage,
			File:                src,
			Filename:            file.Filename,
		}
		updatedPDF, filename, err := s.SignFile(c.Request().Context(), req)
		if err != nil {
			return err
		}

		c.Response().Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		c.Blob(http.StatusOK, "application/pdf", updatedPDF)

		return nil
	}
}
