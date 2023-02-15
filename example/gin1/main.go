package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	helmet "github.com/danielkov/gin-helmet"
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
	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: log.Entry.Logger,
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
	mysqlDBname := os.Getenv("MYSQL_DBNAME")
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlServer := os.Getenv("MYSQL_SERVER")
	mysqlPort := os.Getenv("MYSQL_PORT")

	mysqlConfig := &infrastructure.MySQLConfig{
		Config: &mysql.Config{
			DSN:                       db.GetMySQLDsn(mysqlUsername, mysqlPassword, mysqlServer, mysqlPort, mysqlDBname, "charset=utf8mb4&parseTime=True&loc=Local"),
			DefaultStringSize:         256,   // default size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
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
			ConnMaxLifetime: time.Minute * 5,
			// ConnMaxIdletime sets max idle time(sec)
			ConnMaxIdletime: time.Minute * 5,
			// MaxIdleConns sets idle connection
			MaxIdleConns: 20,
			// MaxOpenConns sets max connection
			MaxOpenConns: 100,
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

	s3Config := infrastructure.GetS3Config(log, s3ID, s3Secret, s3Token, s3Region, s3Endpoint, true)
	sess := infrastructure.NewS3Session(&session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	s3 := s3.New(sess, s3Config)
	awsS3Repository := repository.NewAWSS3Repository(s3, sess)

	// --------------------------------------------------------------
	// Redis
	// --------------------------------------------------------------
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsername := os.Getenv("REDIS_Username")
	redisPassword := os.Getenv("REDIS_Password")

	o := &redis.Options{
		Addr:     redisAddr,
		Username: redisUsername,
		Password: redisPassword,
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
			MaxAge:           12 * time.Hour,
		}))
	router.Use(helmet.Default())
	router.Use(middleware.GinHTTPLogger(log, "request-id", "test"))
	router.GET("/healthcheck", h.GetHealthcheck)
	router.GET("/hello", h.GetHello)
	router.GET("/error_1", h.GetError1)
	router.GET("/error_2", h.GetError2)
	router.GET("/mysql", h.GetMySQL)
	router.GET("/postgres", h.GetPostgres)
	router.GET("/s3", h.GetS3)
	router.GET("/redis", h.GetRedis)
	router.GET("/env", h.GetEnv)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     router,
		ReadTimeout: 5 * time.Second,
	}

	// detect signal
	signal.DetectSignal(func(sig os.Signal) {
		log.Infof("signal detected: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Info("server shutdown...")
		if err = server.Shutdown(ctx); err != nil {
			log.WithField("err", err).Error("server Shutdown")
		}
	}, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, os.Interrupt)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
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
	if err := database.Close(); err != nil {
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
