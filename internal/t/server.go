package t

import (
	"testproject/internal/env"

	"gorm.io/gorm"
)

type Server interface {
	DB() *gorm.DB
	ENV() *env.Env
	HaStart()
	HaStop()
	HaIsRunning() bool
	HaGenerateConfig() error
	HaCheckConfig() error
	HaGetStats() (*[]ProxyStatus, error)
	HaKeepAlive()
}
