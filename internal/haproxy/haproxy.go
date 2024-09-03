package haproxy

import (
	"crypto/rand"
	"fmt"
	"io"
	"os/exec"
	"strconv"
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
	onStopChan       chan bool
	keepAliveEnabled bool
	stdOut           io.Writer
	stdErr           io.Writer
}

type HaproxyInternal struct {
	isRunning   bool
	UUID        string
	Cmd         *exec.Cmd
	isReloading bool
	sync.Mutex
}

func NewHaproxy(s t.Server) *Haproxy {
	stdOut, stderr := NewStdLog(s.DB())
	h := &Haproxy{
		i: HaproxyInternal{
			isRunning: false,
			Cmd:       nil,
		},
		stopChan:   make(chan bool),
		onStopChan: make(chan bool, 1),
		s:          s,
		stdOut:     stdOut,
		stdErr:     stderr,
	}

	return h
}

func (h *Haproxy) Start() error {
	if err := h.start(false); err != nil {
		return err
	}
	// force update state
	h.RunKeepAlive()
	return nil
}

func (h *Haproxy) Reload() error {
	if err := h.start(true); err != nil {
		return err
	}
	// force update state
	h.RunKeepAlive()
	return nil
}

func (h *Haproxy) Stop() error {
	h.i.Lock()
	if !h.i.isRunning {
		h.i.Unlock()
		return nil
	}
	h.shouldBeRunning(false)
	h.i.Unlock()
	h.stopChan <- true
	// wait for stop
	select {
	case <-h.onStopChan:
	case <-time.After(time.Second * 15):
		return fmt.Errorf("timeout while waiting for haproxy to stop")
	}

	// force update state
	h.RunKeepAlive()

	return nil
}

func (h *Haproxy) IsRunning() bool {
	h.i.Lock()
	defer h.i.Unlock()
	return h.i.isRunning
}

func (h *Haproxy) ConfigPath() string {
	return h.s.ENV().WorkDir + "/haproxy/haproxy.cfg"
}

// the locking & unlocking should be done by the caller
func (h *Haproxy) shouldBeRunning(should bool) error {
	tx := h.s.DB()
	var setting m.Setting
	if err := tx.First(&setting).Error; err != nil {
		log.Error().Err(err).Msg("failed to get setting inside start")
		tx.Rollback()
		return err
	}
	setting.ShouldProxyRun = should
	if err := tx.Save(&setting).Error; err != nil {
		log.Error().Err(err).Msg("failed to update setting inside start")
		tx.Rollback()
		return err
	}
	return nil
}

func (h *Haproxy) start(reload bool) error {
	r := make([]byte, 8)
	_, err := rand.Read(r)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%x", r)

	h.GenerateConfig(false)

	if reload {
		log.Info().Msgf("[%s]: haproxy is reloading", name)
	} else {
		log.Info().Msgf("[%s]: haproxy is starting", name)
	}

	h.i.Lock()
	if reload {
		if !h.i.isRunning || h.i.Cmd == nil || h.i.Cmd.Process == nil || h.i.Cmd.Process.Pid < 1 {
			h.i.Unlock()
			h.start(false)
			return nil
		}
	}
	defer h.i.Unlock()
	h.i.isReloading = true
	var tmp *exec.Cmd
	if reload && h.s.ENV().Socket {
		// new process with socket
		tmp = exec.Command("haproxy", "-f", h.ConfigPath(), "-x", "/var/run/haproxy.sock", "-sf", strconv.Itoa(h.i.Cmd.Process.Pid))
	} else {
		// kill the old process if exists
		if h.i.Cmd != nil && h.i.Cmd.Process != nil && h.i.Cmd.Process.Pid >= 1 {
			if err := h.i.Cmd.Process.Kill(); err != nil {
				log.Error().Err(err).Msgf("[%s]: failed to kill old haproxy process", name)
			}
		}
		// new process
		tmp = exec.Command("haproxy", "-f", h.ConfigPath())
	}
	// tmp.Stdout = os.Stdout
	// tmp.Stderr = os.Stdout
	tmp.Stdout = h.stdOut
	tmp.Stderr = h.stdErr

	if err := tmp.Start(); err != nil {
		log.Error().Err(err).
			Msgf("[%s]: failed to start haproxy process", name)
	}

	if tmp.Process == nil || tmp.Process.Pid < 1 {
		log.Error().
			Msgf("[%s]: failed to start haproxy process (unknown pid)", name)
		return fmt.Errorf("[%s]: failed to start haproxy process (unknown pid)", name)
	}

	h.i.UUID = name
	h.i.Cmd = tmp
	h.i.isReloading = false
	h.shouldBeRunning(true)

	// listen for exit
	go func() {
		err := h.i.Cmd.Wait()
		if err != nil {
			log.Error().Err(err).Msgf("[%s]: haproxy exited", name)
		} else {
			log.Info().Msgf("[%s]: haproxy has stopped", name)
		}
	}()

	// listen for stop signal
	go func() {
		for {
			<-h.stopChan
			h.i.Lock()
			if h.i.UUID == name {
				log.Info().Msgf("[%s]: haproxy is being stopped", name)
				if h.i.Cmd != nil && h.i.Cmd.Process != nil {
					h.i.Cmd.Process.Kill()
				}
				h.i.Cmd = nil
				h.onStopChan <- true
			}
			h.i.Unlock()
		}
	}()

	return nil
}
