package bg

import (
	"testproject/internal/app"
	"testproject/internal/m"
	"testproject/internal/t"
	"time"

	"github.com/rs/zerolog/log"
)

func KeepAliveProxy(s t.Server) {
	for {
		time.Sleep(time.Second * 5)
		tx := s.DB()
		var setting m.Setting
		if err := tx.First(&setting).Error; err != nil {
			log.Error().Err(err).Msg("failed to get setting inside keepalive")
			continue
		}

		if app.Proxy.IsRunning() != setting.ShouldProxyRun {
			log.Info().
				Bool("IsRunning", app.Proxy.IsRunning()).
				Bool("ShouldBeRunning", setting.ShouldProxyRun).
				Msg("keepalive proxy state mismatch")
			if !setting.ShouldProxyRun {
				log.Info().Msg("keepalive detected proxy running")
				app.Proxy.Stop()
			} else {
				log.Info().Msg("keepalive detected proxy not running")
				app.Proxy.Start()
			}
		}
	}
}