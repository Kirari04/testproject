package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
)

type CreateAcmeCf struct {
	s      t.Server
	values struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	}
}

func NewCreateAcmeCf(s t.Server) *CreateAcmeCf {
	return &CreateAcmeCf{s: s}
}

func (h *CreateAcmeCf) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	tx := h.s.DB().Begin()
	var setting m.Setting
	if err := tx.Model(&m.Setting{}).First(&setting).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&m.AcmeCloudflareDNSAPIToken{}).Create(&m.AcmeCloudflareDNSAPIToken{
		SettingID: setting.ID,
		Name:      h.values.Name,
		Token:     h.values.Token,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

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
