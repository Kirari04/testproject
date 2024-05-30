package db

import (
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog/log"
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
	})

	log.Info().Msg("Migration database")
	start := time.Now()
	if err = m.Migrate(); err != nil {
		return nil, err
	}
	log.Info().Dur("time_taken", time.Since(start)).Msg("Migration database done")

	return db, nil
}
