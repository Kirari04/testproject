package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type ReloadHandler struct {
	s t.Server
}

func NewReloadHandler(s t.Server) *ReloadHandler {
	return &ReloadHandler{s: s}
}

func (h *ReloadHandler) Route(c echo.Context) error {
	if !h.s.HaIsRunning() {
		return c.String(http.StatusOK, "proxy is not running")
	}
	if err := h.s.HaGenerateConfig(); err != nil {
		return err
	}
	return c.String(http.StatusOK, "ok")
}
