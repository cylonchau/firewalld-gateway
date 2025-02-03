package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/config"
)

// KlogLogger 实现 gorm.Logger 接口
type KlogLogger struct {
	logLevel logger.LogLevel
}

// LogMode 设置日志级别
func (l KlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := l
	newLogger.logLevel = level
	return newLogger
}

// Info 记录 info 级别的日志
func (l KlogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info && klog.V(4).Enabled() {
		klog.Infof(msg, data...)
	}
}

// Warn 记录 warn 级别的日志
func (l KlogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn && klog.V(1).Enabled() {
		klog.Warningf(msg, data...)
	}
}

// Error 记录 error 级别的日志
func (l KlogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error && klog.V(1).Enabled() {
		klog.Errorf(msg, data...)
	}
}

// Trace 记录 trace 级别的日志
func (l KlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		// 错误日志输出在 v1 级别
		if klog.V(1).Enabled() {
			klog.Errorf("Trace Error: %v | SQL: %s | Rows affected: %d | Time: %s", err, sql, rows, elapsed)
		}
	} else {
		// 正常 SQL 日志输出在 v4 级别
		if klog.V(4).Enabled() {
			klog.Infof("Trace Success | SQL: %s | Rows affected: %d | Time: %s", sql, rows, elapsed)
		}
	}
}

var dbConn *sql.DB
var DB *gorm.DB

func InitDB(driver string) error {
	var enconterError error
	newLogger := KlogLogger{logLevel: logger.Info}
	switch driver {
	case "mysql":
		dsn := config.CONFIG.MySQL.User + ":" + config.CONFIG.MySQL.Password + "@tcp(" + config.CONFIG.MySQL.IP + ":" + config.CONFIG.MySQL.Port + ")/" + config.CONFIG.MySQL.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
		if DB, enconterError = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger}); enconterError == nil {
			if dbConn, enconterError = DB.DB(); enconterError == nil {
				dbConn.SetMaxOpenConns(config.CONFIG.MySQL.MaxOpenConnection)
				dbConn.SetMaxIdleConns(config.CONFIG.MySQL.MaxIdleConnection)
				klog.V(4).Infof("Databases stats is %+v", dbConn.Stats())
				return nil
			}
		}
	case "sqlite":
		if DB, enconterError = gorm.Open(sqlite.Open(config.CONFIG.SQLite.File+".db"), &gorm.Config{Logger: newLogger}); enconterError == nil {
			if dbConn, enconterError = DB.DB(); enconterError == nil {
				dbConn.SetMaxOpenConns(config.CONFIG.SQLite.MaxOpenConnection)
				dbConn.SetMaxIdleConns(config.CONFIG.SQLite.MaxIdleConnection)

				klog.V(4).Infof("Databases stats is %+v", dbConn.Stats())
				return nil
			}
		}
	}

	return enconterError
}
