package handler

import (
	"net/http"
	"testproject/internal/app"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type StopHandler struct {
	s t.Server
}

func NewStopHandler(s t.Server) *StopHandler {
	return &StopHandler{s: s}
}

func (h *StopHandler) Route(c echo.Context) error {
	if !app.Proxy.IsRunning() {
		return c.String(http.StatusOK, "proxy is already stopped")
	}
	app.Proxy.Stop()
	return c.String(http.StatusOK, "ok")
}
