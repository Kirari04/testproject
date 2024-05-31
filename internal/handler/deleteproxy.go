package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
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
	if err := tx.Delete(&m.Frontend{}, h.values.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.String(http.StatusOK, "ok")
}
