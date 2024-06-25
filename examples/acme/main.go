package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"github.com/joho/godotenv"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}

func (u *MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func main() {
	// Set Cloudflare API token
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: os.Getenv("USER_EMAIL"),
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// // This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	// config.CADirURL = "http://192.168.99.100:4000/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the DNS-01 challenge provider to Cloudflare

	cfconfig := &cloudflare.Config{
		TTL:                120,
		PropagationTimeout: 2 * time.Minute,
		PollingInterval:    2 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// cfconfig.AuthEmail = os.Getenv("CLOUDFLARE_EMAIL")
	// cfconfig.AuthKey = os.Getenv("CLOUDFLARE_API_KEY")
	cfconfig.AuthToken = os.Getenv("CLOUDFLARE_DNS_API_TOKEN")
	// cfconfig.ZoneToken = os.Getenv("CLOUDFLARE_ZONE_API_TOKEN")
	provider, err := cloudflare.NewDNSProviderConfig(cfconfig)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		log.Fatal(err)
	}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{"test.m141.xyz"},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%s: %s\n", string(certificates.Domain), string(certificates.IssuerCertificate))

	os.WriteFile("cert.pem", certificates.Certificate, 0644)
	os.WriteFile("key.pem", certificates.PrivateKey, 0644)
}
