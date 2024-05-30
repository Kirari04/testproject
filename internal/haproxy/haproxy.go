package haproxy

import (
	"bytes"
	"os/exec"
	"sync"

	"github.com/rs/zerolog/log"
)

type Haproxy struct {
	isRunning bool
	Cmd       *exec.Cmd
	stopChan  chan int
	sync.Mutex
}

func NewHaproxy() *Haproxy {
	return &Haproxy{
		isRunning: false,
		Cmd:       nil,
		stopChan:  make(chan int),
	}
}

func (h *Haproxy) Start() {
	h.Lock()
	if h.isRunning {
		h.Unlock()
		return
	}

	h.Cmd = exec.Command("haproxy", "-f", "./haproxy/haproxy.cfg")
	var stdOut, stdErr bytes.Buffer
	h.Cmd.Stdout = &stdOut
	h.Cmd.Stderr = &stdErr

	log.Info().Msg("haproxy is starting")
	if err := h.Cmd.Start(); err != nil {
		log.Error().Err(err).
			Str("stdout", stdOut.String()).Str("stderr", stdErr.String()).
			Msg("failed to start haproxy")
	}

	log.Info().Msg("haproxy is started")
	h.isRunning = true

	h.Unlock()

	go func() {
		for {
			if h.IsRunning() {
				return
			}
			m, err := stdErr.ReadString('\n')
			if err != nil {
				log.Error().Err(err).Msg("failed to read haproxy stderr")
				h.stopChan <- 1
				continue
			}
			log.Debug().Str("message", m).Msg("haproxy stderr")
		}
	}()
	go func() {
		<-h.stopChan
		h.Lock()
		if h.Cmd != nil && h.Cmd.Process != nil {
			h.Cmd.Process.Kill()
		}
		h.Cmd = nil
		h.isRunning = false
		h.Unlock()
		log.Info().Msg("haproxy is being stopped")
	}()
}

func (h *Haproxy) Stop() {
	h.Lock()
	defer h.Unlock()
	if !h.isRunning {
		return
	}
	h.stopChan <- 1
}

func (h *Haproxy) IsRunning() bool {
	h.Lock()
	defer h.Unlock()
	return h.isRunning
}
