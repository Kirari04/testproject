package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	if err := h.s.HaStart(); err != nil {
		log.Error().Err(err).Msg("failed to start proxy")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start proxy, check crash report for more informations")
	}
	return c.String(http.StatusOK, "ok")
}
