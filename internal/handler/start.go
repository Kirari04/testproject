package handler

import (
	"net/http"
	"testproject/internal/app"

	"github.com/labstack/echo/v4"
)

func Start(c echo.Context) error {
	if app.Proxy.IsRunning() {
		return c.String(http.StatusOK, "already running")
	}
	app.Proxy.Start()
	return c.String(http.StatusOK, "ok")
}
