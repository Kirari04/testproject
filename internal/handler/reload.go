package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ReloadHandler struct {
	s t.Server
}

func NewReloadHandler(s t.Server) *ReloadHandler {
	return &ReloadHandler{s: s}
}

func (h *ReloadHandler) Route(c echo.Context) error {
	if err := h.s.HaReload(); err != nil {
		log.Error().Err(err).Msg("failed to reload proxy")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to reload proxy")
	}
	return c.String(http.StatusOK, "ok")
}
