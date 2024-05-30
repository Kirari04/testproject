package handler

import (
	"net/http"
	"testproject/internal/app"

	"github.com/labstack/echo/v4"
)

func Status(c echo.Context) error {
	if !app.Proxy.IsRunning() {
		return c.String(http.StatusOK, "no")
	}
	return c.String(http.StatusOK, "ok")
}
