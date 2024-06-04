package m

import "time"

type Frontend struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// This is the Port on what the frontend is listening
	Port int `gorm:"column:port" json:"port"`
	// This is the Domain on what the Access Rule will be based on
	Domain string `gorm:"column:domain" json:"domain"`

	// Default Upload Bandwith Limit
	DefBwInLimit uint `gorm:"column:bw_limit" json:"bw_limit"`
	// Default Upload Bandwith Limit Unit
	DefBwInLimitUnit uint `gorm:"column:bw_limit_unit" json:"bw_limit_unit"`
	// Default Upload Bandwith Period in seconds
	DefBwInPeriod uint `gorm:"column:bw_period" json:"bw_period"`
	// Default Download Bandwith Limit
	DefBwOutLimit uint `gorm:"column:bw_out_limit" json:"bw_out_limit"`
	// Default Download Bandwith Limit Unit
	DefBwOutLimitUnit uint `gorm:"column:bw_out_limit_unit" json:"bw_out_limit_unit"`
	// Default Download Bandwith Period in seconds
	DefBwOutPeriod uint `gorm:"column:bw_out_period" json:"bw_out_period"`

	// Default Ratelimit
	DefRateLimit uint `gorm:"column:rate_limit" json:"rate_limit"`
	// Default Ratelimit Period in seconds
	DefRatePeriod uint `gorm:"column:rate_period" json:"rate_period"`
	// Default Hard Ratelimit
	DefHardRateLimit uint `gorm:"column:hard_rate_limit" json:"hard_rate_limit"`
	// Default Hard Ratelimit Period in seconds
	DefHardRatePeriod uint `gorm:"column:hard_rate_period" json:"hard_rate_period"`

	Backends []Backend `json:"backends"`
}
