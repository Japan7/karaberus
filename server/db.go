package server

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var db_instance *gorm.DB = nil

func init_db() {
	if CONFIG.DB.Delete {
		err := os.Remove(CONFIG.DB.File)
		// probably errors don't matter
		if err != nil {
			log.Warn(err.Error())
		}
	}

	fmt.Printf("DB file: %s\n", CONFIG.DB.File)
	db, err := gorm.Open(sqlite.Open(CONFIG.DB.File), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	db_instance = db
}

func GetDB() *gorm.DB {
	return db_instance
}
