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
	db_file := getEnvDefault("DB_FILE", "karaberus.db")
	db_test := getEnvDefault("TEST", "")

	if db_test != "" {
		err := os.Remove(db_file)
		// probably errors don't matter
		if err != nil {
			log.Warn(err.Error())
		}
	}

	fmt.Printf("DB file: %s", db_file)
	db, err := gorm.Open(sqlite.Open(db_file), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	db_instance = db
}

func GetDB() *gorm.DB {
	return db_instance
}
