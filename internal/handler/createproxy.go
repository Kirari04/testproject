package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type CreateProxyHandler struct {
	s      t.Server
	values struct {
		Port     int    `json:"port"`
		Domain   string `json:"domain"`
		Backends []struct {
			Address string `json:"address"`
		}
	}
}

func NewCreateProxyHandler(s t.Server) *CreateProxyHandler {
	return &CreateProxyHandler{s: s}
}

func (h *CreateProxyHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	tx := h.s.DB().Begin()
	frontend := m.Frontend{
		Port:   h.values.Port,
		Domain: h.values.Domain,
	}
	if err := tx.Create(&frontend).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, backend := range h.values.Backends {
		backend := m.Backend{
			Address:    backend.Address,
			FrontendID: frontend.ID,
		}
		if err := tx.Create(&backend).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.String(http.StatusOK, "ok")
}
