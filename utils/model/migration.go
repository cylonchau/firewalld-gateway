package model

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/config"
)

func Migration(driver string) error {
	if driver == "" {
		return errors.New("Unkown database driver")
	}
	var (
		dbInterface   *gorm.DB
		enconterError error
	)
	switch driver {
	case "mysql":
	case "sqlite":
		if dbInterface, enconterError = SQLite(); enconterError == nil {
			dbInterface.Migrator().CurrentDatabase()
			if !dbInterface.Migrator().HasTable(&User{}) {
				dbInterface.Migrator().CreateTable(&User{})
			}
			if !dbInterface.Migrator().HasTable(&Tag{}) {
				dbInterface.Migrator().CreateTable(&Tag{})
			}
			if !dbInterface.Migrator().HasTable(&Host{}) {
				dbInterface.Migrator().CreateTable(&Host{})
			}
			if !dbInterface.Migrator().HasTable(&Template{}) {
				dbInterface.Migrator().CreateTable(&Template{})
			}
			if !dbInterface.Migrator().HasTable(&Port{}) {
				dbInterface.Migrator().CreateTable(&Port{})
			}
			if !dbInterface.Migrator().HasTable(&Rich{}) {
				dbInterface.Migrator().CreateTable(&Rich{})
			}
			if !dbInterface.Migrator().HasTable(&Token{}) {
				dbInterface.Migrator().CreateTable(&Token{})
			}
		}
		return nil
	}
	return enconterError
}

func SQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(config.CONFIG.SQLite.File+".db"), &gorm.Config{})
}

func MySQL() (*gorm.DB, error) {
	return nil, nil
}
