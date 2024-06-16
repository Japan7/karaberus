package server

import (
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db_instance *gorm.DB = nil

func init_db() {
	if CONFIG.DB.Delete {
		err := os.Remove(CONFIG.DB.File)
		// probably errors don't matter
		if err != nil {
			Warn(err.Error())
		}
	}

	gorm_logger := logger.New(
		getLogger(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
		},
	)

	getLogger().Printf("DB file: %s\n", CONFIG.DB.File)
	db, err := gorm.Open(sqlite.Open(CONFIG.DB.File), &gorm.Config{
		Logger: gorm_logger,
	})
	if err != nil {
		panic("Could not connect to the database")
	}

	db_instance = db
}

func GetDB() *gorm.DB {
	return db_instance
}
