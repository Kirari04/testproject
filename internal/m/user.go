package m

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"-"`
}
