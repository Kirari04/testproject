package t

import (
	"testproject/internal/env"

	"gorm.io/gorm"
)

type Server interface {
	DB() *gorm.DB
	ENV() *env.Env
	HaStart() error
	HaStop() error
	HaReload() error
	HaIsRunning() bool
	HaGenerateConfig(reload bool) error
	HaCheckConfig() error
	HaGetStats() (*[]ProxyStatus, error)
	HaStartKeepAlive()
	HaStopKeepAlive()
	HaConfigPath() string
	HaGetCrashReasons() HaproxyCrashReasonsData
}

type HaproxyCrashReasonsData struct {
	HasCrashed bool `json:"has_crashed"`

	AddressInUse    bool   `json:"address_in_use"`
	AddressInUseLog string `json:"address_in_use_log"`

	PermissionDeniedPort    bool   `json:"permission_denied_port"`
	PermissionDeniedPortLog string `json:"permission_denied_port_log"`
}
