package m

import "time"

type Certificate struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Name    string `gorm:"column:name" json:"name"`
	PemPath string `gorm:"column:pem_path" json:"pem_path"`
}
