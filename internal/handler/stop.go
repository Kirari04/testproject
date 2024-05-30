package handler

import (
	"net/http"
	"testproject/internal/app"

	"github.com/labstack/echo/v4"
)

func Stop(c echo.Context) error {
	if !app.Proxy.IsRunning() {
		return c.String(http.StatusOK, "already not running")
	}
	app.Proxy.Stop()
	return c.String(http.StatusOK, "ok")
}
