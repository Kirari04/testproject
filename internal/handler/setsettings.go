package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type SetSettingsHandler struct {
	s      t.Server
	values struct {
		AcmeEmail                 string `json:"acme_email"`
		AcmeCloudflareDNSAPIToken string `json:"acme_cloudflare_dns_api_token"`
	}
}

func NewSetSettingsHandler(s t.Server) *SetSettingsHandler {
	return &SetSettingsHandler{s: s}
}

func (h *SetSettingsHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	tx := h.s.DB().Begin()
	var setting m.Setting
	if err := tx.Model(&m.Setting{}).First(&setting).Error; err != nil {
		tx.Rollback()
		return err
	}

	setting.AcmeEmail = h.values.AcmeEmail
	setting.AcmeCloudflareDNSAPIToken = h.values.AcmeCloudflareDNSAPIToken

	if err := tx.Model(&setting).Save(&setting).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return c.String(http.StatusOK, "ok")
}
