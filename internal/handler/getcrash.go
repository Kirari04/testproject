package handler

import (
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type GetCrashHandler struct {
	s t.Server
}

func NewGetCrashHandler(s t.Server) *GetCrashHandler {
	return &GetCrashHandler{s: s}
}

func (h *GetCrashHandler) Route(c echo.Context) error {
	return c.JSON(http.StatusOK, h.s.HaGetCrashReasons())
}
