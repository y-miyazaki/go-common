package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/example/postgres/entity"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// --------------------------------------------------------------
	// logrus
	// --------------------------------------------------------------
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	logger := infrastructure.NewLogger(logrusLogger)

	// --------------------------------------------------------------
	// logger for gorm
	// --------------------------------------------------------------
	loggerGorm := infrastructure.NewLoggerGorm(&infrastructure.LoggerGormConfig{
		Logger: logger.Entry.Logger,
		GormConfig: &infrastructure.GormConfig{
			// slow query time: 3 sec
			SlowThreshold:             time.Second * 3,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  infrastructure.Info,
		},
	})
	gc := &gorm.Config{
		Logger: loggerGorm,
	}

	// --------------------------------------------------------------
	// MySQL
	// --------------------------------------------------------------
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDBname := os.Getenv("POSTGRES_DBNAME")
	postgresConfig := &infrastructure.PostgresConfig{
		Config: &postgres.Config{
			DSN:                  utils.GetPostgresDsn(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDBname, "sslmode=disable TimeZone=Asia/Tokyo"),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		},
		DBConfig: infrastructure.DBConfig{
			// ConnMaxLifetime sets max life time(sec)
			ConnMaxLifetime: time.Minute * 5,
			// ConnMaxIdletime sets max idle time(sec)
			ConnMaxIdletime: time.Minute * 5,
			// MaxIdleConns sets idle connection
			MaxIdleConns: 20,
			// MaxOpenConns sets max connection
			MaxOpenConns: 100,
		},
	}
	db := infrastructure.NewPostgres(postgresConfig, gc)

	// --------------------------------------------------------------
	// example: Postgres
	// --------------------------------------------------------------
	db.Migrator().CreateTable(&entity.User{})
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	db.Create(user1)

	user2 := &entity.User{}
	db.Take(user2)
	db.Migrator().DropTable(&entity.User{})

	logrusLogger.Infof("name = %s, email = %s", user2.Name, user2.Email)
}
