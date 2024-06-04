package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
)

type GetProxiesStatusHandler struct {
	s t.Server
}

func NewGetProxiesStatusHandler(s t.Server) *GetProxiesStatusHandler {
	return &GetProxiesStatusHandler{s: s}
}

func (h *GetProxiesStatusHandler) Route(c echo.Context) error {
	// Step 1: Fetch the metrics
	url := "http://127.0.0.1:8405/metrics"
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch metrics")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Msg("failed to fetch metrics")
		return errors.New("failed to fetch metrics")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read metrics")
		return err
	}

	// Step 2: Parse the fetched metrics
	parser := expfmt.TextParser{}
	metricFamilies, err := parser.TextToMetricFamilies(bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse metrics")
		return err
	}

	apiRes := make([]t.ProxyStatus, 0)
	// Step 3: Extract the required data
	for name, mf := range metricFamilies {
		if name == "haproxy_server_check_status" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var proxy, server, state string
				for _, label := range labels {
					if label.GetName() == "proxy" {
						proxy = label.GetValue()
					} else if label.GetName() == "server" {
						server = label.GetValue()
					} else if label.GetName() == "state" {
						state = label.GetValue()
					}
				}
				value := metric.GetGauge().GetValue()
				apiRes = append(apiRes, t.ProxyStatus{
					Proxy:  proxy,
					Server: server,
					State:  state,
					Value:  value,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, apiRes)
}
