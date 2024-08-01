package handler

import (
	"errors"
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GetProxyHandler struct {
	s      t.Server
	values struct {
		ID uint `json:"id" query:"id"`
	}
}

func NewGetProxyHandler(s t.Server) *GetProxyHandler {
	return &GetProxyHandler{s: s}
}

func (h *GetProxyHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	tx := h.s.DB()
	var res m.Frontend
	if err := tx.Model(&m.Frontend{}).
		Preload("Backends").
		Preload("Aliases").
		First(&res, h.values.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Proxy not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, res)
}
