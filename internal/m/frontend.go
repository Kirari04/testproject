package m

import "time"

type Frontend struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Port   int    `gorm:"column:port" json:"port"`
	Domain string `gorm:"column:domain" json:"domain"`

	Backends []Backend `json:"backends"`
}
