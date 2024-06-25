package handler

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"testproject/internal/m"
	"testproject/internal/t"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RequestCertificateHandler struct {
	s      t.Server
	values struct {
		Domain string `json:"domain"`
	}
}

func NewRequestCertificateHandler(s t.Server) *RequestCertificateHandler {
	return &RequestCertificateHandler{s: s}
}

func (h *RequestCertificateHandler) Route(c echo.Context) error {
	if err := c.Bind(&h.values); err != nil {
		return err
	}

	if h.values.Domain == "" {
		return c.String(http.StatusBadRequest, "domain is required")
	}

	db := h.s.DB()

	var setting m.Setting
	if err := db.Model(&m.Setting{}).First(&setting).Error; err != nil {
		return err
	}

	certDir := h.s.ENV().WorkDir + "/certs"
	pemName := fmt.Sprintf("%s.pem", uuid.New().String())
	pemPath := certDir + "/" + pemName

	tx := h.s.DB().Begin()

	if err := tx.Create(&m.Certificate{
		Name:    h.values.Domain,
		PemPath: pemName,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add certificate to database: %v", err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	myUser := t.NewAcmeUser(setting.AcmeEmail, privateKey)

	config := lego.NewConfig(myUser)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create lego client: %v", err)
	}

	cfconfig := &cloudflare.Config{
		TTL:                120,
		PropagationTimeout: 2 * time.Minute,
		PollingInterval:    2 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	cfconfig.AuthToken = setting.AcmeCloudflareDNSAPIToken

	provider, err := cloudflare.NewDNSProviderConfig(cfconfig)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create cloudflare dns provider: %v", err)
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set dns01 provider: %v", err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to register: %v", err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{h.values.Domain},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to obtain certificate: %v", err)
	}

	// Destination
	dst, err := os.Create(pemPath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	src := bytes.NewReader(certificates.Certificate)
	if _, err = io.Copy(dst, src); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to copy certificate file: %v", err)
	}

	src2 := bytes.NewReader(certificates.PrivateKey)
	if _, err = io.Copy(dst, src2); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to copy key file: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// reload haproxy
	if err := h.s.HaGenerateConfig(true); err != nil {
		log.Error().Err(err).Msg("Failed to generate proxy config")
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to generate proxy config: Check logs for more information")
	}

	return c.String(http.StatusOK, "ok")
}
