package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type GetHaproxyLogsHandler struct {
	s t.Server
}

func NewGetHaproxyLogsHandler(s t.Server) *GetHaproxyLogsHandler {
	return &GetHaproxyLogsHandler{s: s}
}

func (h *GetHaproxyLogsHandler) Route(c echo.Context) error {
	tx := h.s.DB()
	res := make([]m.HaproxyLog, 0)
	if err := tx.Model(&m.HaproxyLog{}).Order("id DESC").Limit(1000).Find(&res).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
