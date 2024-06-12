package haproxy

import (
	"bytes"
	"sync"
	"testproject/internal/m"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type LogStdBuf struct {
	isTracking bool
	out        bytes.Buffer
	err        bytes.Buffer
	sync.Mutex
}

func (l *LogStdBuf) ReadStdOut() (string, error) {
	return l.out.ReadString('\n')
}

func (l *LogStdBuf) ReadStdErr() (string, error) {
	return l.err.ReadString('\n')
}

func (l *LogStdBuf) Track(db *gorm.DB) {
	l.Lock()
	if l.isTracking {
		l.Unlock()
		return
	}
	l.isTracking = true
	l.Unlock()
	go func() {
		for {
			msg, err := l.ReadStdOut()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if msg != "" {
				if err := db.Create(&m.HaproxyLog{
					Data: msg,
				}).Error; err != nil {
					log.Error().Err(err).Msg("failed to save haproxy log")
				}
			}
		}
	}()
	go func() {
		for {
			msg, err := l.ReadStdErr()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if msg != "" {
				if err := db.Create(&m.HaproxyLog{
					Data: msg,
				}).Error; err != nil {
					log.Error().Err(err).Msg("failed to save haproxy log")
				}
			}
		}
	}()
}

func NewStdLog() (*LogStdBuf, bytes.Buffer, bytes.Buffer) {
	out := bytes.Buffer{}
	err := bytes.Buffer{}
	return &LogStdBuf{
		out: out,
		err: err,
	}, out, err
}
