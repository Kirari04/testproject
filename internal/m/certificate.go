package m

import "time"

var (
	AuthTypes = []string{"cloudflare_dns_api_token"}
)

type Certificate struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	AuthType string `gorm:"column:auth_type" json:"auth_type"`
	AuthID   string `gorm:"column:auth_id" json:"auth_id"`

	Name    string `gorm:"column:name" json:"name"`
	PemPath string `gorm:"column:pem_path" json:"pem_path"`
}
