package haproxy

import (
	"errors"
	"sync"
	"testproject/internal/m"
	"time"

	"github.com/rs/zerolog/log"
)

type HaproxyKeepaliveState struct {
	sync.Mutex
	keepAliveEnabled bool
	// this values should be only used inside RunKeepAlive
	// it's used to detect if haproxy is running or not
	// if i value is higher than 5 no start will be attempted anymore
	invalidRunningStateCounter int
	stopChan                   chan bool
}

func (h *Haproxy) keepAlive() {
	h.RunKeepAlive()
	for {
		time.Sleep(time.Second * 2)
		h.RunKeepAlive()
		if !h.IsKeepAliveEnabled() {
			log.Info().Msg("keepalive stopped")
			h.keepaliveState.stopChan <- true
			break
		}
	}
}

func (h *Haproxy) StartKeepAlive() {
	h.keepaliveState.Lock()
	if h.keepaliveState.keepAliveEnabled {
		h.keepaliveState.Unlock()
		return
	}
	h.keepaliveState.stopChan = make(chan bool)
	h.keepaliveState.keepAliveEnabled = true
	h.keepaliveState.Unlock()
	log.Info().Msg("[keepalive] starting keepalive")
	h.crashReasons.Reset()
	go h.keepAlive()
}

func (h *Haproxy) StopKeepAlive() {
	// stop keepalive
	h.keepaliveState.Lock()
	if !h.keepaliveState.keepAliveEnabled {
		h.keepaliveState.Unlock()
		return
	}
	log.Info().Msg("[keepalive] stopping keepalive")
	h.keepaliveState.keepAliveEnabled = false
	h.keepaliveState.Unlock()
	<-h.keepaliveState.stopChan
}

func (h *Haproxy) RunKeepAlive() {
	// run keepalive
	tx := h.s.DB()
	var setting m.Setting
	if err := tx.First(&setting).Error; err != nil {
		log.Error().Err(err).Msg("[keepalive] failed to get setting")
		return
	}
	if h.IsRunningCheck() != setting.ShouldProxyRun {
		log.Info().
			Bool("IsRunning", h.IsRunning()).
			Bool("IsRunningCheck", h.IsRunningCheck()).
			Bool("ShouldBeRunning", setting.ShouldProxyRun).
			Int("Start Attempts", h.keepaliveState.GetInvalidRunningStateCounter()).
			Msg("[keepalive] proxy state mismatch")
		if !setting.ShouldProxyRun {
			log.Info().Msg("[keepalive] detected proxy running, but expectes is stopped")
			log.Info().Msg("[keepalive] force stopping haproxy")
			h.Stop()
		} else {
			h.keepaliveState.IncrementInvalidRunningStateCounter()
			if h.keepaliveState.GetInvalidRunningStateCounter() > 5 {
				log.Info().Msg("[keepalive] detected proxy failed to start after 5 attempts")
				log.Info().Msg("[keepalive] force stopping haproxy")
				log.Info().Msg("[keepalive] force stopping keepalive")
				h.StopKeepAlive()
				h.Stop()
				return
			}
			log.Info().Msg("[keepalive] detected proxy not running")
			if err := h.Start(); err != nil {
				log.Error().Err(err).Msg("[keepalive] failed to start haproxy")
				return
			}
			log.Info().Msg("[keepalive] successfully started haproxy")
			h.keepaliveState.SetInvalidRunningStateCounter(0)
		}
	}
}

func (h *Haproxy) WaitforIsRunningState() error {
	// we check for 10 seconds if haproxy is running and no crashes happend
	// if it sucessfully started we will reset the invalidRunningStateCounter back to 0
	var crashed bool
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		if !h.IsRunningCheck() {
			crashed = true
			break
		}
	}

	isRunning := h.IsRunningCheck()
	h.i.Lock()
	h.i.isRunning = isRunning
	h.i.Unlock()

	if !crashed {
		h.keepaliveState.SetInvalidRunningStateCounter(0)
		h.crashReasons.Reset()
		return nil
	}
	log.Info().Msg("detected proxy crashed, generating crash report")
	if err := h.GenerateCrashReport(); err != nil {
		log.Error().Err(err).Msg("keepalive failed to generate crash report")
	}

	return errors.New("failed to reach running state")
}

func (k *HaproxyKeepaliveState) GetInvalidRunningStateCounter() int {
	k.Lock()
	defer k.Unlock()
	return k.invalidRunningStateCounter
}

func (k *HaproxyKeepaliveState) IncrementInvalidRunningStateCounter() {
	k.Lock()
	defer k.Unlock()
	k.invalidRunningStateCounter++
}

func (k *HaproxyKeepaliveState) SetInvalidRunningStateCounter(i int) {
	k.Lock()
	defer k.Unlock()
	k.invalidRunningStateCounter = i
}

func (h *Haproxy) IsKeepAliveEnabled() bool {
	h.keepaliveState.Lock()
	defer h.keepaliveState.Unlock()
	return h.keepaliveState.keepAliveEnabled
}
