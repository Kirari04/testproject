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
		AcmeEmail string `json:"acme_email"`
		Acme      struct {
			Cf []struct {
				ID    uint   `json:"id"`
				Name  string `json:"name"`
				Token string `json:"token"`
			} `json:"cf"`
		} `json:"acme"`
	}
}

func NewGetSettingsHandler(s t.Server) *GetSettingsHandler {
	return &GetSettingsHandler{s: s}
}

func (h *GetSettingsHandler) Route(c echo.Context) error {
	tx := h.s.DB()

	var setting m.Setting
	if err := tx.
		Model(&m.Setting{}).
		First(&setting).Error; err != nil {
		return err
	}

	h.res.AcmeEmail = setting.AcmeEmail

	if err := tx.Model(&m.AcmeCloudflareDNSAPIToken{}).
		Find(&h.res.Acme.Cf).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, h.res)
}
