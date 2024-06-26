package haproxy

import (
	"testproject/internal/m"
	"time"

	"github.com/rs/zerolog/log"
)

func (h *Haproxy) KeepAlive() {
	h.keepAliveEnabled = true
	h.RunKeepAlive()
	for {
		time.Sleep(time.Second * 2)
		h.RunKeepAlive()
		if !h.keepAliveEnabled {
			log.Info().Msg("keepalive stopped")
			break
		}
	}
}

func (h *Haproxy) StopKeepAlive() {
	// stop keepalive
	h.keepAliveEnabled = false
}

func (h *Haproxy) RunKeepAlive() {
	// update isRunning
	h.i.Lock()
	if h.i.Cmd == nil || h.i.Cmd.Process == nil || h.i.Cmd.Process.Pid < 1 || h.i.Cmd.ProcessState != nil {
		// log.Debug().Msg("keepalive: haproxy is not running")
		h.i.isRunning = false
	} else {
		// log.Debug().Int("pid", h.i.Cmd.Process.Pid).Msg("keepalive: haproxy is running")
		h.i.isRunning = true
	}
	h.i.Unlock()

	// run keepalive
	tx := h.s.DB()
	var setting m.Setting
	if err := tx.First(&setting).Error; err != nil {
		log.Error().Err(err).Msg("failed to get setting inside keepalive")
		return
	}
	if h.IsRunning() != setting.ShouldProxyRun {
		log.Info().
			Bool("IsRunning", h.IsRunning()).
			Bool("ShouldBeRunning", setting.ShouldProxyRun).
			Msg("keepalive proxy state mismatch")
		if !setting.ShouldProxyRun {
			log.Info().Msg("keepalive detected proxy running")
			h.Stop()
		} else {
			log.Info().Msg("keepalive detected proxy not running")
			h.Start()
		}
	}
}
