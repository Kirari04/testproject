package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type StartHandler struct {
	s t.Server
}

func NewStartHandler(s t.Server) *StartHandler {
	return &StartHandler{s: s}
}

func (h *StartHandler) Route(c echo.Context) error {
	if h.s.HaIsRunning() {
		return c.String(http.StatusOK, "proxy is already running")
	}
	if err := h.s.HaGenerateConfig(); err != nil {
		return err
	}
	return c.String(http.StatusOK, "ok")
}
