package database

import (
	"fmt"
	"time"

	"github.com/gpiecyk/data-warehouse/internal/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSLMode  string
	TimeZone string
}

func (config *Config) getDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.User,
		config.Password,
		config.Dbname,
		config.Port,
		config.SSLMode,
		config.TimeZone,
	)
}

func NewService(config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.getDSN()))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&users.User{})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
