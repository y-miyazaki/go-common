package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/example/mysql/entity"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetPostgres handler
func (h *HTTPHandler) GetPostgres(c *gin.Context) {
	// --------------------------------------------------------------
	// logger for gorm
	// --------------------------------------------------------------
	loggerGorm := logger.NewLoggerGorm(&logger.LoggerGormConfig{
		Logger: h.Logger.Entry.Logger,
		GormConfig: &logger.GormConfig{
			// slow query time: 3 sec
			SlowThreshold:             time.Second * 3,
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
	err := db.Migrator().CreateTable(&entity.User{})
	if err != nil {
		panic("can't create table")
	}
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	_ = db.Create(user1)

	user2 := &entity.User{}
	db.Take(user2)
	err = db.Migrator().DropTable(&entity.User{})
	if err != nil {
		panic("can't drop table")
	}

	h.Logger.Infof("name = %s, email = %s", user2.Name, user2.Email)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
