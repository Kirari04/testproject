package handler

import (
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"
	"testproject/internal/util"

	"github.com/labstack/echo/v4"
)

type CreateProxyHandler struct {
	s      t.Server
	values struct {
		Port   int    `json:"port"`
		Domain string `json:"domain"`

		BwInLimit      uint `json:"bw_in_limit"`
		BwInLimitUnit  uint `json:"bw_in_limit_unit"`
		BwInPeriod     uint `json:"bw_in_period"`
		BwOutLimit     uint `json:"bw_out_limit"`
		BwOutLimitUnit uint `json:"bw_out_limit_unit"`
		BwOutPeriod    uint `json:"bw_out_period"`
		RateLimit      uint `json:"rate_limit"`
		RatePeriod     uint `json:"rate_period"`
		HardRateLimit  uint `json:"hard_rate_limit"`
		HardRatePeriod uint `json:"hard_rate_period"`

		Https       *bool `json:"https"`
		HttpsVerify *bool `json:"https_verify"`

		HttpCheck             *bool   `json:"http_check"`
		HttpCheckMethod       *string `json:"http_check_method"`
		HttpCheckPath         *string `json:"http_check_path"`
		HttpCheckExpectStatus *int    `json:"http_check_expect_status"`
		HttpCheckInterval     *int    `json:"http_check_interval"`
		HttpCheckFailAfter    *int    `json:"http_check_fail_after"`
		HttpCheckRecoverAfter *int    `json:"http_check_recover_after"`
		Backends              []struct {
			Address string `json:"address"`
		}
		Aliases []struct {
			Domain string `json:"domain"`
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
	https := false
	if h.values.Https != nil {
		https = *h.values.Https
	}
	httpsVerify := false
	if h.values.HttpsVerify != nil {
		httpsVerify = *h.values.HttpsVerify
	}

	tx := h.s.DB().Begin()
	frontend := m.Frontend{
		Port:   h.values.Port,
		Domain: h.values.Domain,

		DefBwInLimit:      h.values.BwInLimit,
		DefBwInLimitUnit:  h.values.BwInLimitUnit,
		DefBwInPeriod:     h.values.BwInPeriod,
		DefBwOutLimit:     h.values.BwOutLimit,
		DefBwOutLimitUnit: h.values.BwOutLimitUnit,
		DefBwOutPeriod:    h.values.BwOutPeriod,
		DefRateLimit:      h.values.RateLimit,
		DefRatePeriod:     h.values.RatePeriod,
		DefHardRateLimit:  h.values.HardRateLimit,
		DefHardRatePeriod: h.values.HardRatePeriod,

		HttpCheck:             h.values.HttpCheck,
		HttpCheckMethod:       h.values.HttpCheckMethod,
		HttpCheckPath:         h.values.HttpCheckPath,
		HttpCheckExpectStatus: h.values.HttpCheckExpectStatus,
		HttpCheckInterval:     h.values.HttpCheckInterval,
		HttpCheckFailAfter:    h.values.HttpCheckFailAfter,
		HttpCheckRecoverAfter: h.values.HttpCheckRecoverAfter,
	}
	if err := tx.Create(&frontend).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, backendRaw := range h.values.Backends {
		backend := m.Backend{
			FrontendID: frontend.ID,
			Address:    backendRaw.Address,

			Https:       https,
			HttpsVerify: httpsVerify,
		}
		if err := tx.Create(&backend).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, aliasRaw := range h.values.Aliases {
		alias := m.Alias{
			Domain:     aliasRaw.Domain,
			FrontendID: frontend.ID,
		}
		if err := tx.Create(&alias).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := util.GenerateProxyConfig(h.s); err != nil {
		return err
	}

	return c.String(http.StatusOK, "ok")
}
