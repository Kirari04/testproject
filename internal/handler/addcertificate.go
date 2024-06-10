package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AddCertificateHandler struct {
	s      t.Server
	values struct {
		Name string `form:"name"`
	}
}

func NewAddCertificateHandler(s t.Server) *AddCertificateHandler {
	return &AddCertificateHandler{s: s}
}

func (h *AddCertificateHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	certFile, err := c.FormFile("crt")
	if err != nil {
		log.Warn().Err(err).Msg("no certificate file")
		return echo.NewHTTPError(http.StatusBadRequest, "no certificate file")
	}
	keyFile, err := c.FormFile("key")
	if err != nil {
		log.Warn().Err(err).Msg("no key file")
		return echo.NewHTTPError(http.StatusBadRequest, "no key file")
	}

	tx := h.s.DB().Begin()

	certDir := h.s.ENV().WorkDir + "/certs"
	pemName := fmt.Sprintf("%s.pem", uuid.New().String())
	pemPath := certDir + "/" + pemName

	if err := tx.Create(&m.Certificate{
		Name:    h.values.Name,
		PemPath: pemName,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add certificate to database: %v", err)
	}

	// Destination
	dst, err := os.Create(pemPath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	src, err := certFile.Open()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to open certificate file: %v", err)
	}
	defer src.Close()

	if _, err = io.Copy(dst, src); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to copy certificate file: %v", err)
	}

	src2, err := keyFile.Open()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to open key file: %v", err)
	}
	defer src2.Close()

	if _, err = io.Copy(dst, src2); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to copy key file: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return c.String(http.StatusOK, "ok")
}
