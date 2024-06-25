package t

import (
	"crypto"

	"github.com/go-acme/lego/v4/registration"
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func NewAcmeUser(email string, key crypto.PrivateKey) *AcmeUser {
	return &AcmeUser{
		Email:        email,
		key:          key,
		Registration: nil,
	}
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}

func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
