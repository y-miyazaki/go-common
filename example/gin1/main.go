// Package main demonstrates a complete Gin web application with database connections,
// Redis integration, S3 operations, and various HTTP handlers.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/example/gin1/handler"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/middleware"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"github.com/y-miyazaki/go-common/pkg/signal"
	"github.com/y-miyazaki/go-common/pkg/utils/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	l := &logrus.Logger{}
	// formatter
	l.Formatter = &logrus.JSONFormatter{}
	// out
	l.Out = os.Stdout
	// level
	level, err := logrus.ParseLevel("info")
	if err != nil {
		panic(fmt.Sprintf("level can't set %v", level))
	}
	l.Level = level
	log := logger.NewLogger(l)

	// --------------------------------------------------------------
	// logger for gorm
	// --------------------------------------------------------------
	const slowThresholdSeconds = 3
	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: log.Entry.Logger,
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

	const defaultStringSize = 256
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100
	const corsMaxAgeHours = 12
	const serverReadTimeoutSeconds = 5
	const shutdownTimeoutSeconds = 5
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
	mysqlDB := infrastructure.NewMySQL(mysqlConfig, gc)
	defer closeDB(log, mysqlDB)

	// --------------------------------------------------------------
	// Postgres
	// --------------------------------------------------------------
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDBname := os.Getenv("POSTGRES_DBNAME")
	postgresConfig := &infrastructure.PostgresConfig{
		Config: &postgres.Config{
			DSN:                  db.GetPostgresDsn(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDBname, "sslmode=disable TimeZone=Asia/Tokyo"),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
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
	postgresDB := infrastructure.NewPostgres(postgresConfig, gc)
	defer closeDB(log, postgresDB)
	// --------------------------------------------------------------
	// S3(minio)
	// --------------------------------------------------------------
	s3Region := os.Getenv("S3_REGION")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ID := os.Getenv("S3_ID")
	s3Secret := os.Getenv("S3_SECRET")
	s3Token := os.Getenv("S3_TOKEN")

	s3Config, err := infrastructure.GetAWSS3Config(log, s3ID, s3Secret, s3Token, s3Region, s3Endpoint, true)
	if err != nil {
		panic(err)
	}
	awsS3Repository := repository.NewAWSS3Repository(s3.NewFromConfig(s3Config, func(o *s3.Options) { o.UsePathStyle = true }))

	// --------------------------------------------------------------
	// Redis
	// --------------------------------------------------------------
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsername := os.Getenv("REDIS_Username")
	redisPassword := os.Getenv("REDIS_Password")

	o := &redis.Options{
		Addr:     redisAddr,
		Username: redisUsername,
		Password: redisPassword, // pragma: allowlist secret
	}
	r := infrastructure.NewRedis(o)
	redisRepository := repository.NewRedisRepository(r)
	defer closeRedis(log, r)

	// --------------------------------------------------------------
	// Handler
	// --------------------------------------------------------------
	h := handler.NewHTTPHandler(log, mysqlDB, postgresDB, awsS3Repository, redisRepository)

	router := gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(middleware.GinCors(
		&middleware.GinCorsConfig{
			AllowAllOrigins:  false,
			AllowOrigins:     []string{"https://foo.com"},
			AllowMethods:     []string{"PUT", "PATCH"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           corsMaxAgeHours * time.Hour,
		}))
	router.Use(ginhelmet.Default())
	router.Use(middleware.GinHTTPLogger(log, "request-id", "test"))
	router.GET("/healthcheck", h.HealthCheck)
	router.GET("/hello", h.SayHello)
	router.GET("/error_1", h.HandleError1)
	router.GET("/error_2", h.HandleError2)
	router.GET("/mysql", h.HandleMySQL)
	router.GET("/postgres", h.HandlePostgres)
	router.GET("/s3", h.HandleS3)
	router.GET("/redis", h.HandleRedis)
	router.GET("/env", h.HandleEnv)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     router,
		ReadTimeout: serverReadTimeoutSeconds * time.Second,
	}

	// detect signal
	signal.DetectSignal(func(sig os.Signal) {
		log.Infof("signal detected: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
		defer cancel()
		log.Info("server shutdown...")
		if err = server.Shutdown(ctx); err != nil {
			log.WithField("err", err).Error("server Shutdown")
		}
	}, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, os.Interrupt)

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithField("err", err).Errorf("an error occurred in Server")
	}
}

func closeDB(log *logger.Logger, gormDB *gorm.DB) {
	if gormDB == nil {
		return
	}
	database, err := gormDB.DB()
	if err != nil {
		log.Errorf("can't close db")
		return
	}
	if closeErr := database.Close(); closeErr != nil {
		log.Errorf("can't close db")
	}
}

func closeRedis(log *logger.Logger, r *redis.Client) {
	if r == nil {
		return
	}
	if err := r.Close(); err != nil {
		log.Errorf("can't close redis")
	}
}
