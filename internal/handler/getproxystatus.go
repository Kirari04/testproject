package handler

import (
	"net/http"
	"testproject/internal/t"
	"testproject/internal/util"

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
	apiRes, err := util.GetHaproxyStats(h.s)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch haproxy stats")
		return err
	}

	return c.JSON(http.StatusOK, apiRes)
}
