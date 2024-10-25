package database

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(config *Config) (*gorm.DB, error) {
	l := logger.New(
		log.StandardLogger(),
		logger.Config{
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
			SlowThreshold:             time.Second,
			Colorful:                  false,
		},
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: l})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
