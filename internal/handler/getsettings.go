package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type GetSettingsHandler struct {
	s   t.Server
	res struct {
		AcmeEmail                 string `json:"acme_email"`
		AcmeCloudflareDNSAPIToken string `json:"acme_cloudflare_dns_api_token"`
	}
}

func NewGetSettingsHandler(s t.Server) *GetSettingsHandler {
	return &GetSettingsHandler{s: s}
}

func (h *GetSettingsHandler) Route(c echo.Context) error {
	tx := h.s.DB()
	if err := tx.Model(&m.Setting{}).
		First(&h.res).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, h.res)
}
