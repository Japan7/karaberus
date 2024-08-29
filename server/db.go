package server

import (
	"context"
	"errors"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db_instance *gorm.DB = nil

func init_db() {
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

	gorm_config := &gorm.Config{
		Logger: gorm_logger,
	}

	if CONFIG.DB.Driver == "sqlite" {
		getLogger().Printf("DB file: %s\n", CONFIG.DB.File)
		db, err := gorm.Open(sqlite.Open(CONFIG.DB.File), gorm_config)
		if err != nil {
			panic("Could not connect to the database")
		}
		sqlite_db, err := db.DB()
		if err != nil {
			panic("Could not get underlying sqlite db")
		}
		sqlite_db.SetMaxOpenConns(1)

		db_instance = db
	} else if CONFIG.DB.Driver == "postgres" {
		getLogger().Printf("Postgres DSN: %s\n", CONFIG.DB.DSN)
		db, err := gorm.Open(postgres.Open(CONFIG.DB.DSN))
		if err != nil {
			panic("Could not connect to the database")
		}
		db_instance = db
	} else {
		panic("unknown db driver " + CONFIG.DB.Driver)
	}
}

func GetDB(ctx context.Context) *gorm.DB {
	if db_instance == nil {
		init_db()
		init_model(db_instance)
	}

	return db_instance.WithContext(ctx)
}

func DBErrToHumaErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return huma.Error404NotFound("record not found")
	}
	return err
}
