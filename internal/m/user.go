package m

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
}
