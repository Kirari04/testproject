package haproxy

import (
	"strings"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/rs/zerolog/log"
)

func (h *Haproxy) GenerateCrashReport() error {
	log.Warn().Msg("Generating crash report")
	tx := h.s.DB()
	var haproxyLogs []m.HaproxyLog
	if err := tx.Model(&m.HaproxyLog{}).Order("id desc").Limit(10).Find(&haproxyLogs).Error; err != nil {
		log.Error().Err(err).Msg("failed to get haproxy logs")
		return err
	}
	h.crashReasons.Lock()
	h.crashReasons.HasCrashed = true
	var foundPermissionDeniedPort bool
	var foundAddressInUse bool
	for _, haproxyLog := range haproxyLogs {
		if strings.Contains(haproxyLog.Data, "cannot bind socket (Permission denied) for") && !foundPermissionDeniedPort {
			foundPermissionDeniedPort = true
			h.crashReasons.PermissionDeniedPort = true
			h.crashReasons.PermissionDeniedPortLog = haproxyLog.Data
		}
		if strings.Contains(haproxyLog.Data, "cannot bind socket (Address already in use) for") && !foundAddressInUse {
			foundAddressInUse = true
			h.crashReasons.AddressInUse = true
			h.crashReasons.AddressInUseLog = haproxyLog.Data
		}
	}
	h.crashReasons.Unlock()
	return nil
}

func (h *Haproxy) GetCrashReasons() t.HaproxyCrashReasonsData {
	h.crashReasons.Lock()
	defer h.crashReasons.Unlock()
	return t.HaproxyCrashReasonsData{
		HasCrashed:              h.crashReasons.HasCrashed,
		AddressInUse:            h.crashReasons.AddressInUse,
		AddressInUseLog:         h.crashReasons.AddressInUseLog,
		PermissionDeniedPort:    h.crashReasons.PermissionDeniedPort,
		PermissionDeniedPortLog: h.crashReasons.PermissionDeniedPortLog,
	}
}

func (c *HaproxyCrashReasons) Reset() {
	c.Lock()
	defer c.Unlock()
	c.HasCrashed = false
	c.AddressInUse = false
	c.AddressInUseLog = ""
	c.PermissionDeniedPort = false
	c.PermissionDeniedPortLog = ""
}
