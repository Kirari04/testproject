package cmd

import (
	"fmt"
	"log"
	"testproject/internal/util"
	"time"

	"github.com/urfave/cli/v2"
)

func certgen(c *cli.Context) error {
	host := "example.com"
	validFrom := ""
	validFor := 365 * 24 * time.Hour
	isCA := false
	rsaBits := 2048
	ecdsaCurve := ""
	ed25519Key := false

	cert, key, err := util.GenerateSelfSignedCert(host, validFrom, validFor, isCA, rsaBits, ecdsaCurve, ed25519Key)
	if err != nil {
		log.Fatalf("Error generating certificate: %v", err)
	}

	fmt.Println("Certificate:\n", cert)
	fmt.Println("Key:\n", key)

	return nil
}
