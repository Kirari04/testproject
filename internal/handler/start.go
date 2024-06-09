package handler

import (
	"net/http"
	"testproject/internal/app"
	"testproject/internal/t"
	"testproject/internal/util"

	"github.com/labstack/echo/v4"
)

type StartHandler struct {
	s t.Server
}

func NewStartHandler(s t.Server) *StartHandler {
	return &StartHandler{s: s}
}

func (h *StartHandler) Route(c echo.Context) error {
	if app.Proxy.IsRunning() {
		return c.String(http.StatusOK, "proxy is already running")
	}
	if err := util.GenerateProxyConfig(h.s); err != nil {
		return err
	}
	return c.String(http.StatusOK, "ok")
}
