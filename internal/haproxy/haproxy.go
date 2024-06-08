package haproxy

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
	"testproject/internal/m"
	"testproject/internal/t"
	"time"

	"github.com/rs/zerolog/log"
)

type Haproxy struct {
	i        HaproxyInternal
	s        t.Server
	stopChan chan bool
}

type HaproxyInternal struct {
	isRunning bool
	Cmd       *exec.Cmd

	sync.Mutex
}

func NewHaproxy(s t.Server) *Haproxy {
	return &Haproxy{
		i: HaproxyInternal{
			isRunning: false,
			Cmd:       nil,
		},
		stopChan: make(chan bool),
		s:        s,
	}
}

func (h *Haproxy) Start() {
	h.i.Lock()
	if h.i.isRunning {
		h.i.Unlock()
		return
	}

	h.i.Cmd = exec.Command("haproxy", "-f", "./haproxy/haproxy.cfg")
	var stdOut, stdErr bytes.Buffer
	h.i.Cmd.Stdout = &stdOut
	h.i.Cmd.Stderr = &stdErr

	log.Info().Msg("haproxy is starting")
	if err := h.i.Cmd.Start(); err != nil {
		log.Error().Err(err).
			Str("stdout", stdOut.String()).Str("stderr", stdErr.String()).
			Msg("failed to start haproxy")
	}

	log.Info().Msg("haproxy is started")
	tx := h.s.DB()
	h.i.isRunning = true

	var setting m.Setting
	if err := tx.First(&setting).Error; err != nil {
		log.Error().Err(err).Msg("failed to get setting inside start")
		tx.Rollback()
		h.i.Unlock()
		return
	}
	setting.ShouldProxyRun = true
	if err := tx.Save(&setting).Error; err != nil {
		log.Error().Err(err).Msg("failed to update setting inside start")
		tx.Rollback()
		h.i.Unlock()
		return
	}

	h.i.Unlock()

	// output logs
	go func() {
		for {
			if !h.IsRunning() {
				log.Debug().Msg("stop tracking haproxy stderr")
				return
			}
			m, err := stdErr.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(time.Second * 1)
					continue
				}
				log.Error().Err(err).Msg("failed to read haproxy stderr")
				return
			}
			if m != "" {
				log.Debug().Msgf("haproxy: %s", m)
			}
		}
	}()
	go func() {
		for {
			if !h.IsRunning() {
				log.Debug().Msg("stop tracking haproxy stdout")
				return
			}
			m, err := stdOut.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(time.Second * 1)
					continue
				}
				log.Error().Err(err).Msg("failed to read haproxy stdout")
				return
			}
			if m != "" {
				log.Debug().Msgf("haproxy: %s", m)
			}
		}
	}()
	// listen for exit
	go func() {
		err := h.i.Cmd.Wait()
		if err != nil {
			log.Error().Err(err).Msg("haproxy exited")
		}
		h.stopChan <- true
	}()
	// listen for stop signal
	go func() {
		for {
			<-h.stopChan
			h.i.Lock()
			if h.i.Cmd != nil && h.i.Cmd.Process != nil {
				h.i.Cmd.Process.Kill()
			}
			h.i.Cmd = nil

			tx := h.s.DB()
			h.i.isRunning = false

			var setting m.Setting
			if err := tx.First(&setting).Error; err != nil {
				log.Error().Err(err).Msg("failed to get setting inside stop")
				tx.Rollback()
				h.i.Unlock()
				return
			}
			setting.ShouldProxyRun = false
			if err := tx.Save(&setting).Error; err != nil {
				log.Error().Err(err).Msg("failed to update setting inside stop")
				tx.Rollback()
				h.i.Unlock()
				return
			}

			h.i.Unlock()
			log.Info().Msg("haproxy is being stopped")

		}
	}()
}

func (h *Haproxy) Stop() {
	h.i.Lock()
	defer h.i.Unlock()
	if !h.i.isRunning {
		return
	}
	h.stopChan <- true
}

func (h *Haproxy) IsRunning() bool {
	h.i.Lock()
	defer h.i.Unlock()
	return h.i.isRunning
}
