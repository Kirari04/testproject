package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type GetCertificatesHandler struct {
	s t.Server
}

func NewGetCertificatesHandler(s t.Server) *GetCertificatesHandler {
	return &GetCertificatesHandler{s: s}
}

func (h *GetCertificatesHandler) Route(c echo.Context) error {
	tx := h.s.DB()
	res := make([]m.Certificate, 0)
	if err := tx.Model(&m.Certificate{}).Find(&res).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
