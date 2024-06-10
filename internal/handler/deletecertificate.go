package handler

import (
	"fmt"
	"net/http"
	"os"
	"testproject/internal/m"
	"testproject/internal/t"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type DeleteCertificateHandler struct {
	s      t.Server
	values struct {
		ID uint `json:"id"`
	}
}

func NewDeleteCertificateHandler(s t.Server) *DeleteCertificateHandler {
	return &DeleteCertificateHandler{s: s}
}

func (h *DeleteCertificateHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	tx := h.s.DB().Begin()

	var cert m.Certificate
	if err := tx.Model(&cert).First(&cert, h.values.ID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get certificate from database: %v", err)
	}

	if err := tx.Model(&m.Certificate{}).Delete(&cert).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add certificate to database: %v", err)
	}
	certDir := h.s.ENV().WorkDir + "/certs"
	pemPath := certDir + "/" + cert.PemPath
	if err := os.Remove(pemPath); err != nil {
		log.Warn().Err(err).Msg("failed to remove certificate file")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return c.String(http.StatusOK, "ok")
}
