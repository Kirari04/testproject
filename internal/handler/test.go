package handler

import (
	"testproject/internal/util"

	"github.com/labstack/echo/v4"
)

func Test(c echo.Context) error {
	if err := util.TestHaproxyConfig(); err != nil {
		return err
	}
	return c.String(200, "ok")
}
