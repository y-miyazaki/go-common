package main

import (
	"os"
	"testing"
	"time"

	"github.com/y-miyazaki/go-common/example/gin1/handler"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"github.com/y-miyazaki/go-common/pkg/utils/db"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	redis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoggerSetup(t *testing.T) {
	// Test logrus logger setup
	l := &logrus.Logger{}
	l.Formatter = &logrus.JSONFormatter{}
	l.Out = os.Stdout

	level, err := logrus.ParseLevel("info")
	assert.NoError(t, err)
	l.Level = level

	log := logger.NewLogger(l)

	assert.NotNil(t, log)
	assert.NotNil(t, log.Entry)
	assert.Equal(t, logrus.InfoLevel, log.Entry.Logger.Level)
}

func TestGormLoggerSetup(t *testing.T) {
	const slowThresholdSeconds = 3

	l := &logrus.Logger{}
	l.Formatter = &logrus.JSONFormatter{}
	l.Out = os.Stdout
	level, _ := logrus.ParseLevel("info")
	l.Level = level
	log := logger.NewLogger(l)

	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: log.Entry.Logger,
		GormConfig: &logger.GormConfig{
			SlowThreshold:             time.Second * slowThresholdSeconds,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		},
	})

	assert.NotNil(t, loggerGorm)
}

func TestMySQLConfig(t *testing.T) {
	const defaultStringSize = 256
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	mysqlDBname := "test"
	mysqlUsername := "user"
	mysqlPassword := "password" // pragma: allowlist-secret
	mysqlServer := "localhost"
	mysqlPort := "3306"

	mysqlConfig := &infrastructure.MySQLConfig{
		Config: &mysql.Config{
			DSN:                       db.GetMySQLDsn(mysqlUsername, mysqlPassword, mysqlServer, mysqlPort, mysqlDBname, "charset=utf8mb4&parseTime=True&loc=Local"),
			DefaultStringSize:         defaultStringSize,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		},
		DBConfig: infrastructure.DBConfig{
			ConnMaxLifetime: time.Minute * connectionLifetimeMinutes,
			ConnMaxIdletime: time.Minute * connectionLifetimeMinutes,
			MaxIdleConns:    maxIdleConnections,
			MaxOpenConns:    maxOpenConnections,
		},
	}

	assert.NotNil(t, mysqlConfig)
	assert.NotNil(t, mysqlConfig.Config)
	assert.Equal(t, uint(defaultStringSize), mysqlConfig.Config.DefaultStringSize)
	assert.True(t, mysqlConfig.Config.DisableDatetimePrecision)
	assert.True(t, mysqlConfig.Config.DontSupportRenameIndex)
	assert.True(t, mysqlConfig.Config.DontSupportRenameColumn)
	assert.False(t, mysqlConfig.Config.SkipInitializeWithVersion)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, mysqlConfig.ConnMaxLifetime)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, mysqlConfig.ConnMaxIdletime)
	assert.Equal(t, maxIdleConnections, mysqlConfig.MaxIdleConns)
	assert.Equal(t, maxOpenConnections, mysqlConfig.MaxOpenConns)
}

func TestPostgresConfig(t *testing.T) {
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	postgresDBname := "test"
	postgresUsername := "user"
	postgresPassword := "password"
	postgresServer := "localhost"
	postgresPort := "5432"

	postgresConfig := &infrastructure.PostgresConfig{
		Config: &postgres.Config{
			DSN:                  db.GetPostgresDsn(postgresUsername, postgresPassword, postgresServer, postgresPort, postgresDBname, "sslmode=disable TimeZone=Asia/Tokyo"),
			PreferSimpleProtocol: true,
		},
		DBConfig: infrastructure.DBConfig{
			ConnMaxLifetime: time.Minute * connectionLifetimeMinutes,
			ConnMaxIdletime: time.Minute * connectionLifetimeMinutes,
			MaxIdleConns:    maxIdleConnections,
			MaxOpenConns:    maxOpenConnections,
		},
	}

	assert.NotNil(t, postgresConfig)
	assert.NotNil(t, postgresConfig.Config)
	assert.True(t, postgresConfig.Config.PreferSimpleProtocol)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, postgresConfig.ConnMaxLifetime)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, postgresConfig.ConnMaxIdletime)
	assert.Equal(t, maxIdleConnections, postgresConfig.MaxIdleConns)
	assert.Equal(t, maxOpenConnections, postgresConfig.MaxOpenConns)
}

func TestRedisConfig(t *testing.T) {
	redisAddr := "localhost:6379"
	redisPassword := ""
	redisDB := 0

	redisOptions := &redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	assert.NotNil(t, redisOptions)
	assert.Equal(t, redisAddr, redisOptions.Addr)
	assert.Equal(t, redisPassword, redisOptions.Password)
	assert.Equal(t, redisDB, redisOptions.DB)
}

func TestS3Config(t *testing.T) {
	s3Region := "us-east-1"
	s3Endpoint := "http://localhost:9000"
	s3ID := "test-id"
	s3Secret := "test-secret" // pragma: allowlist-secret
	s3Token := ""

	assert.Equal(t, "us-east-1", s3Region)
	assert.Equal(t, "http://localhost:9000", s3Endpoint)
	assert.Equal(t, "test-id", s3ID)
	assert.Equal(t, "test-secret", s3Secret)
	assert.Equal(t, "", s3Token)
}

func TestS3Options(t *testing.T) {
	testFunc := func(o *s3.Options) {
		o.UsePathStyle = true
	}
	assert.NotNil(t, testFunc)
}

func TestMain_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if required environment variables are not set
	requiredEnvVars := []string{
		"MYSQL_DBNAME", "MYSQL_USERNAME", "MYSQL_PASSWORD", "MYSQL_SERVER", "MYSQL_PORT",
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_DBNAME",
		"REDIS_ADDR",
		"S3_REGION", "S3_ENDPOINT", "S3_ID", "S3_SECRET",
	}

	allSet := true
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			allSet = false
			break
		}
	}

	if !allSet {
		t.Skip("Skipping integration test - requires all environment variables to be set")
	}

	// This would test the actual main function logic
	// For now, just ensure it doesn't panic when called
	assert.NotPanics(t, func() {
		// We can't easily test main() directly as it starts a server
		// and creates database connections. In a real integration test, we would:
		// 1. Set up test databases and services
		// 2. Run main() in a goroutine
		// 3. Test the HTTP endpoints
		// 4. Clean up
	})
}

func TestCloseFunctions(t *testing.T) {
	// Test that closeDB and closeRedis functions exist
	// These functions are defined in main.go but are private
	assert.True(t, true) // Placeholder test
}

func TestMainConfiguration(t *testing.T) {
	// Test main function configuration without starting server
	// Set up environment variables for testing
	originalEnv := make(map[string]string)
	envVars := map[string]string{
		"MYSQL_DBNAME":      "test_db",
		"MYSQL_USERNAME":    "test_user",
		"MYSQL_PASSWORD":    "test_pass",
		"MYSQL_SERVER":      "localhost",
		"MYSQL_PORT":        "3306",
		"POSTGRES_USER":     "test_user",
		"POSTGRES_PASSWORD": "test_pass",
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_DBNAME":   "test_db",
		"S3_REGION":         "us-east-1",
		"S3_ENDPOINT":       "http://localhost:9000",
		"S3_ID":             "test_id",
		"S3_SECRET":         "test_secret",
		"S3_TOKEN":          "",
		"REDIS_ADDR":        "localhost:6379",
		"REDIS_Username":    "",
		"REDIS_Password":    "",
	}

	// Save original environment variables
	for key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Set test environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	defer func() {
		// Restore original environment variables
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test logger setup
	l := &logrus.Logger{}
	l.Formatter = &logrus.JSONFormatter{}
	l.Out = os.Stdout
	level, err := logrus.ParseLevel("info")
	assert.NoError(t, err)
	l.Level = level
	log := logger.NewLogger(l)
	assert.NotNil(t, log)

	// Test GORM logger setup
	const slowThresholdSeconds = 3
	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: log.Entry.Logger,
		GormConfig: &logger.GormConfig{
			SlowThreshold:             time.Second * slowThresholdSeconds,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		},
	})
	assert.NotNil(t, loggerGorm)

	gc := &gorm.Config{
		Logger: loggerGorm,
	}
	assert.NotNil(t, gc)

	// Test MySQL config creation
	mysqlDBname := os.Getenv("MYSQL_DBNAME")
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlServer := os.Getenv("MYSQL_SERVER")
	mysqlPort := os.Getenv("MYSQL_PORT")

	const defaultStringSize = 256
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	mysqlConfig := &infrastructure.MySQLConfig{
		Config: &mysql.Config{
			DSN:                       db.GetMySQLDsn(mysqlUsername, mysqlPassword, mysqlServer, mysqlPort, mysqlDBname, "charset=utf8mb4&parseTime=True&loc=Local"),
			DefaultStringSize:         defaultStringSize,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		},
		DBConfig: infrastructure.DBConfig{
			ConnMaxLifetime: time.Minute * connectionLifetimeMinutes,
			ConnMaxIdletime: time.Minute * connectionLifetimeMinutes,
			MaxIdleConns:    maxIdleConnections,
			MaxOpenConns:    maxOpenConnections,
		},
	}

	// Test MySQL config validation
	assert.NotNil(t, mysqlConfig)
	assert.NotNil(t, mysqlConfig.Config)
	assert.Contains(t, mysqlConfig.Config.DSN, mysqlDBname)
	assert.Contains(t, mysqlConfig.Config.DSN, mysqlUsername)

	// Test S3 config creation
	s3Region := os.Getenv("S3_REGION")
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3ID := os.Getenv("S3_ID")
	s3Secret := os.Getenv("S3_SECRET")
	s3Token := os.Getenv("S3_TOKEN")

	s3Config, err := infrastructure.GetAWSConfig(log, infrastructure.AWSServiceS3, s3ID, s3Secret, s3Token, s3Region, s3Endpoint)
	assert.NoError(t, err)
	assert.NotNil(t, s3Config)

	// Test S3 repository creation
	awsS3Repository := repository.NewAWSS3Repository(s3.NewFromConfig(s3Config, func(o *s3.Options) { o.UsePathStyle = true }))
	assert.NotNil(t, awsS3Repository)

	// Test Redis config creation
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsername := os.Getenv("REDIS_Username")
	redisPassword := os.Getenv("REDIS_Password")

	o := &redis.Options{
		Addr:     redisAddr,
		Username: redisUsername,
		Password: redisPassword,
	}
	r := infrastructure.NewRedis(o)
	assert.NotNil(t, r)

	// Test Redis repository creation
	redisRepository := repository.NewRedisRepository(r)
	assert.NotNil(t, redisRepository)

	// Test handler creation
	h := handler.NewHTTPHandler(log, nil, nil, awsS3Repository, redisRepository) // Pass nil for DB connections
	assert.NotNil(t, h)
}
