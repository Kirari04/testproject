package m

import "time"

type Setting struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	ShouldProxyRun bool `gorm:"column:should_proxy_run" json:"should_proxy_run"`

	AcmeEmail                 string `gorm:"column:acme_email" json:"acme_email"`
	AcmeCloudflareDNSAPIToken string `gorm:"column:acme_cloudflare_dns_api_token" json:"acme_cloudflare_dns_api_token"`
}
