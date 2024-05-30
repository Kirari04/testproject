package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type GetProxiesHandler struct {
	s t.Server
}

func NewGetProxiesHandler(s t.Server) *GetProxiesHandler {
	return &GetProxiesHandler{s: s}
}

func (h *GetProxiesHandler) Route(c echo.Context) error {
	tx := h.s.DB().Begin()
	res := make([]m.Frontend, 0)
	if err := tx.Model(&m.Frontend{}).Preload("Backends").Find(&res).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
