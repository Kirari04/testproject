package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
	s t.Server
}

func NewStatusHandler(s t.Server) *StatusHandler {
	return &StatusHandler{s: s}
}

func (h *StatusHandler) Route(c echo.Context) error {
	if !h.s.HaIsRunning() {
		return c.String(http.StatusOK, "no")
	}
	return c.String(http.StatusOK, "ok")
}
