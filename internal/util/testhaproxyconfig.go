package util

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

func TestHaproxyConfig() error {
	cmd := exec.Command("haproxy", "-c", "-V", "-f", "haproxy/haproxy.cfg")
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	cmd.Run()
	cmd.Wait()

	if !strings.HasPrefix(strings.TrimSpace(stdOut.String()), "Configuration file is valid") {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("[%s] - {%s}", stdOut.String(), stdErr.String()))
	}

	return nil
}
