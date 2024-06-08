package m

import "time"

type Backend struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Address     string `gorm:"column:address" json:"address"`
	Https       bool   `gorm:"column:https" json:"https"`
	HttpsVerify bool   `gorm:"column:https_verify" json:"https_verify"`

	FrontendID uint     `gorm:"index,column:frontend_id" json:"-"`
	Frontend   Frontend `json:"-"`
}
