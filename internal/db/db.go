package db

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./data/db.sqlite3"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				type Frontend struct {
					ID        uint      `gorm:"primaryKey;column:id"`
					CreatedAt time.Time `gorm:"column:created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at"`

					Port   int    `gorm:"column:port"`
					Domain string `gorm:"column:domain"`
				}
				return tx.Migrator().CreateTable(&Frontend{})
			},
			Rollback: func(tx *gorm.DB) error {
				type Frontend struct {
					ID        uint      `gorm:"primaryKey;column:id"`
					CreatedAt time.Time `gorm:"column:created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at"`

					Port   int    `gorm:"column:port"`
					Domain string `gorm:"column:domain"`
				}
				return db.Migrator().DropTable(&Frontend{})
			},
		},
		{
			ID: "2",
			Migrate: func(tx *gorm.DB) error {
				type Backend struct {
					ID        uint      `gorm:"primaryKey;column:id"`
					CreatedAt time.Time `gorm:"column:created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at"`

					Address string `gorm:"column:address"`

					FrontendID uint `gorm:"index,column:frontend_id"`
				}
				return tx.Migrator().CreateTable(&Backend{})
			},
			Rollback: func(tx *gorm.DB) error {
				type Backend struct {
					ID        uint      `gorm:"primaryKey;column:id"`
					CreatedAt time.Time `gorm:"column:created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at"`

					Address string `gorm:"column:address"`

					FrontendID uint `gorm:"index,column:frontend_id"`
				}
				return db.Migrator().DropTable(&Backend{})
			},
		},
		{
			ID: "3",
			Migrate: func(tx *gorm.DB) error {
				type Setting struct {
					ID        uint      `gorm:"primaryKey;column:id" json:"id"`
					CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

					ShouldProxyRun bool `gorm:"column:should_proxy_run" json:"should_proxy_run"`
				}
				return tx.Migrator().CreateTable(&Setting{})
			},
			Rollback: func(tx *gorm.DB) error {
				type Setting struct {
					ID        uint      `gorm:"primaryKey;column:id" json:"id"`
					CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

					ShouldProxyRun bool `gorm:"column:should_proxy_run" json:"should_proxy_run"`
				}
				return db.Migrator().DropTable(&Setting{})
			},
		},
		{
			ID: "4",
			Migrate: func(tx *gorm.DB) error {
				type Setting struct {
					ID        uint      `gorm:"primaryKey;column:id" json:"id"`
					CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

					ShouldProxyRun bool `gorm:"column:should_proxy_run" json:"should_proxy_run"`
				}
				return tx.Create(&Setting{
					ShouldProxyRun: false,
				}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				type Setting struct {
					ID        uint      `gorm:"primaryKey;column:id" json:"id"`
					CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

					ShouldProxyRun bool `gorm:"column:should_proxy_run" json:"should_proxy_run"`
				}
				return tx.Where("1 = 1").Delete(&Setting{}).Error
			},
		},
		{
			ID: "5",
			Migrate: func(tx *gorm.DB) error {
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
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefBwInLimit"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefBwInLimitUnit"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefBwInPeriod"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefBwOutLimit"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefBwOutLimitUnit"); err != nil {
					return err
				}
				return tx.Migrator().AddColumn(&Frontend{}, "DefBwOutPeriod")
			},
			Rollback: func(tx *gorm.DB) error {
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
				}
				if err := db.Migrator().DropColumn(&Frontend{}, "DefBwInLimit"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&Frontend{}, "DefBwInLimitUnit"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&Frontend{}, "DefBwInPeriod"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&Frontend{}, "DefBwOutLimit"); err != nil {
					return err
				}
				if err := db.Migrator().DropColumn(&Frontend{}, "DefBwOutLimitUnit"); err != nil {
					return err
				}
				return db.Migrator().DropColumn(&Frontend{}, "DefBwOutPeriod")
			},
		},
		{
			ID: "6",
			Migrate: func(tx *gorm.DB) error {
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
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefRateLimit"); err != nil {
					return err
				}

				return tx.Migrator().AddColumn(&Frontend{}, "DefRatePeriod")
			},
			Rollback: func(tx *gorm.DB) error {
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
				}
				if err := tx.Migrator().DropColumn(&Frontend{}, "DefRateLimit"); err != nil {
					return err
				}

				return tx.Migrator().DropColumn(&Frontend{}, "DefRatePeriod")
			},
		},
		{
			ID: "7",
			Migrate: func(tx *gorm.DB) error {
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
				}
				if err := tx.Migrator().AddColumn(&Frontend{}, "DefHardRateLimit"); err != nil {
					return err
				}

				return tx.Migrator().AddColumn(&Frontend{}, "DefHardRatePeriod")
			},
			Rollback: func(tx *gorm.DB) error {
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
				}
				if err := tx.Migrator().DropColumn(&Frontend{}, "DefHardRateLimit"); err != nil {
					return err
				}

				return tx.Migrator().DropColumn(&Frontend{}, "DefHardRatePeriod")
			},
		},
		{
			ID: "8",
			Migrate: func(tx *gorm.DB) error {
				type Alias struct {
					ID        uint      `gorm:"primaryKey;column:id" json:"id"`
					CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

					Domain string `gorm:"column:domain" json:"domain"`

					FrontendID uint `gorm:"index,column:frontend_id" json:"-"`
				}
				return tx.Migrator().CreateTable(&Alias{})
			},
			Rollback: func(tx *gorm.DB) error {
				type Alias struct {
					ID        uint      `gorm:"primaryKey;column:id" json:"id"`
					CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
					UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

					Domain string `gorm:"column:domain" json:"domain"`

					FrontendID uint `gorm:"index,column:frontend_id" json:"-"`
				}
				return db.Migrator().DropTable(&Alias{})
			},
		},
	})

	log.Info().Msg("Migration database")
	start := time.Now()
	if err = m.Migrate(); err != nil {
		return nil, err
	}
	log.Info().Dur("time_taken", time.Since(start)).Msg("Migration database done")

	return db, nil
}
