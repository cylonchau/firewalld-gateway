package model

import (
	"database/sql"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/config"
)

var dbConn *sql.DB
var DB *gorm.DB

func InitDB(driver string) error {
	var enconterError error
	switch driver {
	case "sqlite":
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: true,
			},
		)
		if DB, enconterError = gorm.Open(sqlite.Open(config.CONFIG.SQLite.File+".db"), &gorm.Config{Logger: newLogger}); enconterError == nil {
			if dbConn, enconterError = DB.DB(); enconterError == nil {
				dbConn.SetMaxOpenConns(config.CONFIG.SQLite.MaxOpenConnection)
				dbConn.SetMaxIdleConns(config.CONFIG.SQLite.MaxIdleConnection)
				klog.V(4).Infof("Databases stats is %+v", dbConn.Stats())
				//DB.Logger.LogMode(logger.Info)
				return nil
			}
		}
	}

	return enconterError
}
