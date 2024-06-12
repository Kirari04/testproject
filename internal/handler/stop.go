package handler

import (
	"net/http"
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
	if !h.s.HaIsRunning() {
		return c.String(http.StatusOK, "proxy is already stopped")
	}
	h.s.HaStop()
	return c.String(http.StatusOK, "ok")
}
