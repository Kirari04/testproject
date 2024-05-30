package handler

import (
	"net/http"
	"testproject/internal/t"
	"testproject/internal/util"

	"github.com/labstack/echo/v4"
)

type ApplyHandler struct {
	s t.Server
}

func NewApplyHandler(s t.Server) *ApplyHandler {
	return &ApplyHandler{s: s}
}

func (h *ApplyHandler) Route(c echo.Context) error {
	if err := util.GenerateProxyConfig(h.s); err != nil {
		return err
	}
	return c.String(http.StatusOK, "ok")
}
