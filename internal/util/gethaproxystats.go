package util

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testproject/internal/app"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
)

func GetHaproxyStats(s t.Server) (*[]t.ProxyStatus, error) {
	if !app.Proxy.IsRunning() {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "haproxy is not running")
	}

	// Step 1: Fetch the metrics
	url := "http://127.0.0.1:8405/metrics"
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch metrics")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Msg("failed to fetch metrics")
		return nil, errors.New("failed to fetch metrics")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read metrics")
		return nil, err
	}

	// Step 2: Parse the fetched metrics
	parser := expfmt.TextParser{}
	metricFamilies, err := parser.TextToMetricFamilies(bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse metrics")
		return nil, err
	}

	data := map[int]map[int]*t.ProxyStatusServer{}
	dataStats := map[int]*t.ProxyStatusStats{}

	// Step 3: Extract the required data
	for name, mf := range metricFamilies {
		// parse status
		if name == "haproxy_server_check_status" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var proxy, server, state string
				for _, label := range labels {
					if label.GetName() == "server" {
						server = label.GetValue()
					} else if label.GetName() == "state" {
						state = label.GetValue()
					} else if label.GetName() == "proxy" {
						proxy = label.GetValue()
					}
				}
				value := metric.GetGauge().GetValue()

				if !strings.Contains(proxy, "_") {
					continue
				}

				// parse frontend id from proxy name
				prxPlits := strings.Split(proxy, "_")
				if len(prxPlits) != 2 {

					return nil, errors.New("invalid proxy name")
				}
				frontendID, err := strconv.Atoi(prxPlits[1])
				if err != nil {
					log.Error().Err(err).Str("proxy", proxy).Msg("invalid frontend id")
					return nil, err
				}

				// parse server(backend)id and frontendId
				srvPlits := strings.Split(server, "_")
				if len(srvPlits) != 3 {
					log.Error().Str("server", server).Msg("invalid server name")
					return nil, errors.New("invalid server name")
				}
				serverID, err := strconv.Atoi(srvPlits[2])
				if err != nil {
					log.Error().Err(err).Str("server", server).Msg("invalid server id")
					return nil, err
				}

				// init frontend stats struct
				if _, ok := dataStats[frontendID]; !ok {
					dataStats[frontendID] = &t.ProxyStatusStats{
						FrontendId: uint(frontendID),
					}
				}

				// init proxy status server struct
				if _, ok := data[frontendID]; !ok {
					var dbFrontend m.Frontend
					if err := s.DB().
						Model(&m.Frontend{}).
						First(&dbFrontend, frontendID).Error; err != nil {
						log.Error().Err(err).Int("frontend_id", frontendID).Msg("failed to fetch frontend")
						return nil, err
					}
					data[frontendID] = map[int]*t.ProxyStatusServer{}
				}
				// init proxy status server struct
				if _, ok := data[frontendID][serverID]; !ok {
					var dbServer m.Backend
					if err := s.DB().
						Model(&m.Backend{}).
						First(&dbServer, serverID).Error; err != nil {
						log.Error().Err(err).Msg("failed to fetch backend")
						return nil, err
					}
					data[frontendID][serverID] = &t.ProxyStatusServer{
						ServerId: uint(serverID),
						Address:  dbServer.Address,
					}
				}

				// add status
				switch state {
				case "HANA":
					data[frontendID][serverID].HANA = value
				case "SOCKERR":
					data[frontendID][serverID].SOCKERR = value
				case "L4OK":
					data[frontendID][serverID].L4OK = value
				case "L4TOUT":
					data[frontendID][serverID].L4TOUT = value
				case "L4CON":
					data[frontendID][serverID].L4CON = value
				case "L6OK":
					data[frontendID][serverID].L6OK = value
				case "L6TOUT":
					data[frontendID][serverID].L6TOUT = value
				case "L6RSP":
					data[frontendID][serverID].L6RSP = value
				case "L7TOUT":
					data[frontendID][serverID].L7TOUT = value
				case "L7RSP":
					data[frontendID][serverID].L7RSP = value
				case "L7OK":
					data[frontendID][serverID].L7OK = value
				case "L7OKC":
					data[frontendID][serverID].L7OKC = value
				case "L7STS":
					data[frontendID][serverID].L7STS = value
				case "PROCERR":
					data[frontendID][serverID].PROCERR = value
				case "PROCTOUT":
					data[frontendID][serverID].PROCTOUT = value
				case "PROCOK":
					data[frontendID][serverID].PROCOK = value
				default:
					log.Error().Str("state", state).Msg("unknown state from prometheus metrics")
				}
			}
		}
		// parse bytes in and out
		if name == "haproxy_backend_bytes_in_total" || name == "haproxy_backend_bytes_out_total" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var proxy string
				for _, label := range labels {
					if label.GetName() == "proxy" {
						proxy = label.GetValue()
					}
				}
				value := metric.GetCounter().GetValue()

				if !strings.Contains(proxy, "_") {
					continue
				}

				// parse frontend id from proxy name
				prxPlits := strings.Split(proxy, "_")
				if len(prxPlits) != 2 {

					return nil, errors.New("invalid proxy name")
				}
				frontendID, err := strconv.Atoi(prxPlits[1])
				if err != nil {
					log.Error().Err(err).Str("proxy", proxy).Msg("invalid frontend id")
					return nil, err
				}

				// init frontend stats struct
				if _, ok := dataStats[frontendID]; !ok {
					dataStats[frontendID] = &t.ProxyStatusStats{
						FrontendId: uint(frontendID),
					}
				}

				if name == "haproxy_backend_bytes_in_total" {
					dataStats[frontendID].BytesInTotal = &value
				} else if name == "haproxy_backend_bytes_out_total" {
					dataStats[frontendID].BytesOutTotal = &value
				} else {
					log.Warn().Str("name", name).Msg("unknown name")
				}
			}
		}
		// parse requests total
		if name == "haproxy_backend_http_requests_total" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var proxy string
				for _, label := range labels {
					if label.GetName() == "proxy" {
						proxy = label.GetValue()
					}
				}
				value := metric.GetCounter().GetValue()

				if !strings.Contains(proxy, "_") {

					continue
				}

				// parse frontend id from proxy name
				prxPlits := strings.Split(proxy, "_")
				if len(prxPlits) != 2 {

					return nil, errors.New("invalid proxy name")
				}
				frontendID, err := strconv.Atoi(prxPlits[1])
				if err != nil {
					log.Error().Err(err).Str("proxy", proxy).Msg("invalid frontend id")
					return nil, err
				}

				// init frontend stats struct
				if _, ok := dataStats[frontendID]; !ok {
					dataStats[frontendID] = &t.ProxyStatusStats{
						FrontendId: uint(frontendID),
					}
				}
				dataStats[frontendID].RequestsTotal = &value
			}
		}
		// parse responses total
		if name == "haproxy_backend_http_responses_total" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var proxy, code string
				for _, label := range labels {
					if label.GetName() == "proxy" {
						proxy = label.GetValue()
					} else if label.GetName() == "code" {
						code = label.GetValue()
					}
				}
				value := metric.GetCounter().GetValue()

				if !strings.Contains(proxy, "_") {

					continue
				}

				// parse frontend id from proxy name
				prxPlits := strings.Split(proxy, "_")
				if len(prxPlits) != 2 {

					return nil, errors.New("invalid proxy name")
				}
				frontendID, err := strconv.Atoi(prxPlits[1])
				if err != nil {
					log.Error().Err(err).Str("proxy", proxy).Msg("invalid frontend id")
					return nil, err
				}

				// init frontend stats struct
				if _, ok := dataStats[frontendID]; !ok {
					dataStats[frontendID] = &t.ProxyStatusStats{
						FrontendId: uint(frontendID),
					}
				}
				switch code {
				case "1xx":
					dataStats[frontendID].ResponsesTotal1xx = &value
				case "2xx":
					dataStats[frontendID].ResponsesTotal2xx = &value
				case "3xx":
					dataStats[frontendID].ResponsesTotal3xx = &value
				case "4xx":
					dataStats[frontendID].ResponsesTotal4xx = &value
				case "5xx":
					dataStats[frontendID].ResponsesTotal5xx = &value
				default:
					dataStats[frontendID].ResponsesTotalOther = &value
				}
			}
		}
	}

	// convert map to array
	apiRes := make([]t.ProxyStatus, 0)
	for frontendId, frontendMap := range data {
		proxyStatus := t.ProxyStatus{
			FrontendId: uint(frontendId),
			Servers:    make([]t.ProxyStatusServer, 0),

			BytesInTotal:        dataStats[frontendId].BytesInTotal,
			BytesOutTotal:       dataStats[frontendId].BytesOutTotal,
			RequestsTotal:       dataStats[frontendId].RequestsTotal,
			ResponsesTotal1xx:   dataStats[frontendId].ResponsesTotal1xx,
			ResponsesTotal2xx:   dataStats[frontendId].ResponsesTotal2xx,
			ResponsesTotal3xx:   dataStats[frontendId].ResponsesTotal3xx,
			ResponsesTotal4xx:   dataStats[frontendId].ResponsesTotal4xx,
			ResponsesTotal5xx:   dataStats[frontendId].ResponsesTotal5xx,
			ResponsesTotalOther: dataStats[frontendId].ResponsesTotalOther,
		}
		for _, serverStatus := range frontendMap {
			proxyStatus.Servers = append(proxyStatus.Servers, *serverStatus)
		}

		apiRes = append(apiRes, proxyStatus)
	}

	return &apiRes, nil
}
