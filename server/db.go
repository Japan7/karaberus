package server

import (
	"context"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

var db_instance *gorm.DB = nil

func init_db() {
	if CONFIG.DB.Delete {
		err := os.Remove(CONFIG.DB.File)
		// probably errors don't matter
		if err != nil {
			logger.Warn(err.Error())
		}
	}

	log.Printf("DB file: %s\n", CONFIG.DB.File)
	db, err := gorm.Open(sqlite.Open(CONFIG.DB.File), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	db_instance = db
}

func GetDB(ctx context.Context) *gorm.DB {
	return db_instance.WithContext(ctx)
}
