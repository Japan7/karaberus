package server

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db_instance *gorm.DB = nil

func init_db() {
	db_file := getEnvDefault("DB_FILE", "karaberus.db")
	db, err := gorm.Open(sqlite.Open(db_file), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	db_instance = db
}

func GetDB() *gorm.DB {
	return db_instance
}
