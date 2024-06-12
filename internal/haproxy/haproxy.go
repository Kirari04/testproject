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
	i                HaproxyInternal
	s                t.Server
	stopChan         chan bool
	keepAliveEnabled bool
}

type HaproxyInternal struct {
	isRunning bool
	Cmd       *exec.Cmd
	stdErr    bytes.Buffer
	stdOut    bytes.Buffer

	sync.Mutex
}

func NewHaproxy(s t.Server) *Haproxy {
	h := &Haproxy{
		i: HaproxyInternal{
			isRunning: false,
			Cmd:       nil,
			stdErr:    bytes.Buffer{},
			stdOut:    bytes.Buffer{},
		},
		stopChan: make(chan bool),
		s:        s,
	}
	// output logs
	go func() {
		tx := s.DB()
		for {
			msg, err := h.i.stdErr.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(time.Second * 1)
					continue
				}
				log.Error().Err(err).Msg("failed to read haproxy stderr")
				return
			}
			if msg != "" {
				if err := tx.Create(&m.HaproxyLog{
					Data: msg,
				}).Error; err != nil {
					log.Error().Err(err).Msg("failed to save haproxy log")
				}
			}
		}
	}()
	go func() {
		tx := s.DB()
		for {
			msg, err := h.i.stdOut.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(time.Second * 1)
					continue
				}
				log.Error().Err(err).Msg("failed to read haproxy stdout")
				return
			}
			if msg != "" {
				if err := tx.Create(&m.HaproxyLog{
					Data: msg,
				}).Error; err != nil {
					log.Error().Err(err).Msg("failed to save haproxy log")
				}
			}
		}
	}()
	return h
}

func (h *Haproxy) Start() {
	h.i.Lock()
	if h.i.isRunning {
		h.i.Unlock()
		return
	}

	h.i.Cmd = exec.Command("haproxy", "-f", "./haproxy/haproxy.cfg")
	h.i.Cmd.Stdout = &h.i.stdOut
	h.i.Cmd.Stderr = &h.i.stdErr

	log.Info().Msg("haproxy is starting")
	if err := h.i.Cmd.Start(); err != nil {
		log.Error().Err(err).
			Str("stdout", h.i.stdOut.String()).Str("stderr", h.i.stdErr.String()).
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
