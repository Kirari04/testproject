package m

import (
	"time"
)

type Setting struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	ShouldProxyRun bool `gorm:"column:should_proxy_run" json:"should_proxy_run"`

	AcmeEmail                  string                      `gorm:"column:acme_email" json:"acme_email"`
	AcmeCloudflareDNSAPITokens []AcmeCloudflareDNSAPIToken `json:"acme_cloudflare_dns_api_tokens"`
}

type AcmeCloudflareDNSAPIToken struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	SettingID uint    `gorm:"column:setting_id" json:"setting_id"`
	Setting   Setting `json:"setting"`

	Name  string `gorm:"column:name" json:"name"`
	Token string `gorm:"column:token" json:"token"`
}
