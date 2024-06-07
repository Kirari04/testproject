package m

import "time"

type Alias struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Domain string `gorm:"column:domain" json:"domain"`

	FrontendID uint     `gorm:"index,column:frontend_id" json:"-"`
	Frontend   Frontend `json:"-"`
}
