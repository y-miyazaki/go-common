// Package main demonstrates MySQL database connection example.
package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/example/mysql/entity"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/utils/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	const slowThresholdSeconds = 3
	const defaultStringSize = 256
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	// --------------------------------------------------------------
	// logrus
	// --------------------------------------------------------------
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	loggerNew := logger.NewLogger(logrusLogger)

	// --------------------------------------------------------------
	// logger for gorm
	// --------------------------------------------------------------
	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: loggerNew.Entry.Logger,
		GormConfig: &logger.GormConfig{
			// slow query time: 3 sec
			SlowThreshold:             time.Second * slowThresholdSeconds,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		},
	})
	gc := &gorm.Config{
		Logger: loggerGorm,
	}

	// --------------------------------------------------------------
	// MySQL
	// --------------------------------------------------------------
	mysqlDBname := os.Getenv("MYSQL_DBNAME")
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlServer := os.Getenv("MYSQL_SERVER")
	mysqlPort := os.Getenv("MYSQL_PORT")

	mysqlConfig := &infrastructure.MySQLConfig{
		Config: &mysql.Config{
			DSN:                       db.GetMySQLDsn(mysqlUsername, mysqlPassword, mysqlServer, mysqlPort, mysqlDBname, "charset=utf8mb4&parseTime=True&loc=Local"),
			DefaultStringSize:         defaultStringSize, // default size for string fields
			DisableDatetimePrecision:  true,              // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,              // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,              // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false,             // auto configure based on currently MySQL version
		},
		DBConfig: infrastructure.DBConfig{
			// ConnMaxLifetime sets max life time(sec)
			ConnMaxLifetime: time.Minute * connectionLifetimeMinutes,
			// ConnMaxIdletime sets max idle time(sec)
			ConnMaxIdletime: time.Minute * connectionLifetimeMinutes,
			// MaxIdleConns sets idle connection
			MaxIdleConns: maxIdleConnections,
			// MaxOpenConns sets max connection
			MaxOpenConns: maxOpenConnections,
		},
	}
	database := infrastructure.NewMySQL(mysqlConfig, gc)

	// --------------------------------------------------------------
	// example: MySQL
	// --------------------------------------------------------------
	err := database.Migrator().CreateTable(&entity.User{})
	if err != nil {
		panic("can't create table")
	}
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	database.Create(user1)

	user2 := &entity.User{}
	database.Take(user2)
	err = database.Migrator().DropTable(&entity.User{})
	if err != nil {
		panic("can't drop table")
	}

	loggerNew.Infof("name = %s, email = %s", user2.Name, user2.Email)
}
