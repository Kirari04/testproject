package haproxy

import (
	"testproject/internal/m"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type LogBuf struct {
	db        *gorm.DB
	WriteChan chan []byte
}

func (l *LogBuf) Write(p []byte) (n int, err error) {
	l.WriteChan <- p
	return len(p), nil
}
func (l *LogBuf) Track() {
	buffer := make([]byte, 1024*1024)
	bi := 0
	for {
		b, ok := <-l.WriteChan
		if !ok {
			break
		}
		for i := 0; i < len(b); i++ {
			if b[i] == '\n' {
				if err := l.db.Create(&m.HaproxyLog{
					Data: string(buffer[:bi]),
				}).Error; err != nil {
					log.Error().Err(err).Msg("failed to save haproxy log")
					log.Debug().Msgf("haproxy log: %s", string(buffer))
				}

				// reset buffer
				buffer = make([]byte, 1024*1024)
				bi = 0
				continue
			}
			// append char to buffer
			buffer[bi] = b[i]
			bi++
		}
	}
}

func NewStdLog(db *gorm.DB) (*LogBuf, *LogBuf) {
	out := LogBuf{
		db:        db,
		WriteChan: make(chan []byte, 1024),
	}
	err := LogBuf{
		db:        db,
		WriteChan: make(chan []byte, 1024),
	}
	go out.Track()
	go err.Track()

	return &out, &err
}
