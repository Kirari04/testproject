package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type DeleteProxyHandler struct {
	s      t.Server
	values struct {
		ID uint `json:"id"`
	}
}

func NewDeleteProxyHandler(s t.Server) *DeleteProxyHandler {
	return &DeleteProxyHandler{s: s}
}

func (h *DeleteProxyHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	tx := h.s.DB().Begin()
	if err := tx.
		Select("Backends", "Aliases").
		Delete(&m.Frontend{ID: h.values.ID}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// reload haproxy
	if err := h.s.HaGenerateConfig(true); err != nil {
		log.Error().Err(err).Msg("Failed to generate proxy config")
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to generate proxy config: Check logs for more information")
	}

	return c.String(http.StatusOK, "ok")
}
