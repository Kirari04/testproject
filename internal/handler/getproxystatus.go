package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type GetProxiesStatusHandler struct {
	s t.Server
}

func NewGetProxiesStatusHandler(s t.Server) *GetProxiesStatusHandler {
	return &GetProxiesStatusHandler{s: s}
}

func (h *GetProxiesStatusHandler) Route(c echo.Context) error {
	apiRes, err := h.s.HaGetStats()
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch haproxy stats")
		return c.JSON(http.StatusOK, []t.ProxyStatus{})
	}

	return c.JSON(http.StatusOK, apiRes)
}
