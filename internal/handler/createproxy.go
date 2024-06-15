package handler

import (
	"fmt"
	"net/http"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CreateProxyHandler struct {
	s      t.Server
	values struct {
		Port   int    `json:"port"`
		Https  *bool  `json:"https"`
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

		BackendHttps       *bool `json:"backend_https"`
		BackendHttpsVerify *bool `json:"backend_https_verify"`

		HttpCheck             *bool   `json:"http_check"`
		HttpCheckMethod       *string `json:"http_check_method"`
		HttpCheckPath         *string `json:"http_check_path"`
		HttpCheckExpectStatus *int    `json:"http_check_expect_status"`
		HttpCheckInterval     *int    `json:"http_check_interval"`
		HttpCheckFailAfter    *int    `json:"http_check_fail_after"`
		HttpCheckRecoverAfter *int    `json:"http_check_recover_after"`

		RequestBodyLimit     uint `json:"request_body_limit"`
		RequestBodyLimitUnit uint `json:"request_body_limit_unit"`

		Backends []struct {
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
	backendHttps := false
	if h.values.BackendHttps != nil {
		backendHttps = *h.values.BackendHttps
	}
	backendHttpsVerify := false
	if h.values.BackendHttpsVerify != nil {
		backendHttpsVerify = *h.values.BackendHttpsVerify
	}

	https := false
	if h.values.Https != nil {
		https = *h.values.Https
	}
	if https {
		// check if other frontends with the same port are listening on https
		var frontends []m.Frontend
		if err := h.s.DB().Where(&m.Frontend{Port: h.values.Port}).Find(&frontends).Error; err != nil {
			return fmt.Errorf("failed to get frontends from database: %w", err)
		}
		for _, frontend := range frontends {
			if !frontend.Https {
				return c.String(http.StatusBadRequest, "This port is already in use with a non-https frontend")
			}
		}
	}

	tx := h.s.DB().Begin()
	frontend := m.Frontend{
		Port:   h.values.Port,
		Https:  https,
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

		DefRequestBodyLimit:     h.values.RequestBodyLimit,
		DefRequestBodyLimitUnit: h.values.RequestBodyLimitUnit,

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

			Https:       backendHttps,
			HttpsVerify: backendHttpsVerify,
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

	if err := h.s.HaGenerateConfig(true); err != nil {
		log.Error().Err(err).Msg("Failed to generate proxy config")
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to generate proxy config: Check logs for more information")
	}

	return c.String(http.StatusOK, "ok")
}
