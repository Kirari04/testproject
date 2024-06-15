package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	if err := h.s.HaStop(); err != nil {
		log.Error().Err(err).Msg("failed to stop haproxy")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to stop haproxy")
	}
	return c.String(http.StatusOK, "ok")
}
