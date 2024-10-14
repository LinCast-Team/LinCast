package database

import (
	"fmt"
	"time"

	"lincast/models"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(dbPort int, dbHost, dbUser, dbPassword, dbName string) (*gorm.DB, error) {
	l := logger.New(
		log.StandardLogger(),
		logger.Config{
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
			SlowThreshold:             time.Second,
			Colorful:                  false,
		},
	)

	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{Logger: l})
	if err != nil {
		return nil, err
	}

	migrate(db)

	return db, nil
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Podcast{},
		&models.Episode{},
		&models.PlaybackInfo{},
		&models.QueueEpisode{},
		&models.EpisodeProgress{},
	)
	if err != nil {
		log.WithError(errorx.EnsureStackTrace(err)).Panic("error when executing automigration")
	}
}
