package util

import (
	"bytes"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testproject/internal/m"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func TestHaproxyConfig(db *gorm.DB) error {
	cmd := exec.Command("haproxy", "-c", "-V", "-f", "haproxy/haproxy.cfg")
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	cmd.Run()
	cmd.Wait()
	isValid := false
	// save logs
	for {
		msg, err := stdErr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error().Err(err).Msg("failed to read haproxy stderr")
			break
		}
		if msg != "" {
			if err := db.Create(&m.HaproxyLog{
				Data: msg,
			}).Error; err != nil {
				log.Error().Err(err).Msg("failed to save haproxy log")
			}
		}
		if strings.HasPrefix(strings.TrimSpace(msg), "Configuration file is valid") {
			isValid = true
		}
	}
	for {
		msg, err := stdOut.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error().Err(err).Msg("failed to read haproxy stdOut")
			break
		}
		if msg != "" {
			if err := db.Create(&m.HaproxyLog{
				Data: msg,
			}).Error; err != nil {
				log.Error().Err(err).Msg("failed to save haproxy log")
			}
		}
		if strings.HasPrefix(strings.TrimSpace(msg), "Configuration file is valid") {
			isValid = true
		}
	}

	if !isValid {
		return echo.NewHTTPError(http.StatusInternalServerError, "Configuration file is invalid")
	}
	return nil
}
