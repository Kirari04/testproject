package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

func Test(c echo.Context) error {
	// haproxy -c -f /usr/local/etc/haproxy/haproxy.cfg
	cmd := exec.Command("haproxy", "-c", "-f", "/usr/local/etc/haproxy/haproxy.cfg")
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	cmd.Run()
	cmd.Wait()

	if !strings.HasPrefix(strings.TrimSpace(stdOut.String()), "Configuration file is valid") {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("[%s] - {%s}", stdOut.String(), stdErr.String()))
	}

	return c.String(http.StatusOK, fmt.Sprintf("[%s] - {%s}", stdOut.String(), stdErr.String()))
}
