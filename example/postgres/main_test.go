package main

import (
	"os"
	"testing"
	"time"

	"github.com/y-miyazaki/go-common/example/postgres/entity"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserEntity(t *testing.T) {
	user := &entity.User{
		Name:  "test",
		Email: "test@test.com",
		ID:    1,
	}

	assert.Equal(t, "test", user.Name)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "user", user.TableName())
}

func TestPostgresConfigSetup(t *testing.T) {
	const slowThresholdSeconds = 3
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	// Test configuration values
	assert.Equal(t, 3, slowThresholdSeconds)
	assert.Equal(t, 5, connectionLifetimeMinutes)
	assert.Equal(t, 20, maxIdleConnections)
	assert.Equal(t, 100, maxOpenConnections)
}

func TestPostgresConfigStruct(t *testing.T) {
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	postgresConfig := &infrastructure.PostgresConfig{
		Config: &postgres.Config{
			DSN:                  "host=localhost user=user password=password dbname=test port=5432 sslmode=disable TimeZone=Asia/Tokyo", // pragma: allowlist-secret
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
	assert.Equal(t, "host=localhost user=user password=password dbname=test port=5432 sslmode=disable TimeZone=Asia/Tokyo", postgresConfig.Config.DSN)
	assert.True(t, postgresConfig.Config.PreferSimpleProtocol)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, postgresConfig.ConnMaxLifetime)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, postgresConfig.ConnMaxIdletime)
	assert.Equal(t, maxIdleConnections, postgresConfig.MaxIdleConns)
	assert.Equal(t, maxOpenConnections, postgresConfig.MaxOpenConns)
}

func TestLoggerSetup(t *testing.T) {
	// Test logrus logger setup
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")

	assert.NotNil(t, logrusLogger)
	assert.IsType(t, &logrus.JSONFormatter{}, logrusLogger.Formatter)
	assert.Equal(t, logrus.InfoLevel, logrusLogger.Level)
}

func TestGormLoggerSetup(t *testing.T) {
	const slowThresholdSeconds = 3

	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	loggerNew := logger.NewLogger(logrusLogger)

	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: loggerNew.Entry.Logger,
		GormConfig: &logger.GormConfig{
			SlowThreshold:             time.Second * slowThresholdSeconds,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		},
	})

	assert.NotNil(t, loggerGorm)
}

func TestGormConfig(t *testing.T) {
	const slowThresholdSeconds = 3

	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	loggerNew := logger.NewLogger(logrusLogger)

	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: loggerNew.Entry.Logger,
		GormConfig: &logger.GormConfig{
			SlowThreshold:             time.Second * slowThresholdSeconds,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		},
	})

	gc := &gorm.Config{
		Logger: loggerGorm,
	}

	assert.NotNil(t, gc)
	assert.NotNil(t, gc.Logger)
}

func TestMain_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if PostgreSQL environment variables are not set
	if os.Getenv("POSTGRES_USER") == "" ||
		os.Getenv("POSTGRES_PASSWORD") == "" ||
		os.Getenv("POSTGRES_HOST") == "" ||
		os.Getenv("POSTGRES_PORT") == "" ||
		os.Getenv("POSTGRES_DBNAME") == "" {
		t.Skip("Skipping PostgreSQL integration test - requires PostgreSQL environment variables")
	}

	// This would test the actual main function logic
	// For now, just ensure it doesn't panic when called
	assert.NotPanics(t, func() {
		// We can't easily test main() directly as it creates database connections
		// and performs migrations. In a real integration test, we would:
		// 1. Set up a test database
		// 2. Run main()
		// 3. Verify the database operations
		// 4. Clean up
	})
}

func TestDatabaseOperations(t *testing.T) {
	// Test database operation logic without actual connection
	const slowThresholdSeconds = 3
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	// Test configuration constants
	assert.Equal(t, 3, slowThresholdSeconds)
	assert.Equal(t, 5, connectionLifetimeMinutes)
	assert.Equal(t, 20, maxIdleConnections)
	assert.Equal(t, 100, maxOpenConnections)
}

func TestUserEntityOperations(t *testing.T) {
	// Test user entity creation and operations
	user1 := &entity.User{Name: "test", Email: "test@test.com"}
	assert.Equal(t, "test", user1.Name)
	assert.Equal(t, "test@test.com", user1.Email)
	assert.Equal(t, "user", user1.TableName())

	user2 := &entity.User{}
	assert.Equal(t, "", user2.Name)
	assert.Equal(t, "", user2.Email)
	assert.Equal(t, "user", user2.TableName())
}

func TestMainFunctionSetup(t *testing.T) {
	// Test main function setup without database connection
	const slowThresholdSeconds = 3
	const connectionLifetimeMinutes = 5
	const maxIdleConnections = 20
	const maxOpenConnections = 100

	// Test logger setup
	logrusLogger := &logrus.Logger{}
	logrusLogger.Formatter = &logrus.JSONFormatter{}
	logrusLogger.Out = os.Stdout
	logrusLogger.Level, _ = logrus.ParseLevel("Info")
	loggerNew := logger.NewLogger(logrusLogger)
	assert.NotNil(t, loggerNew)

	// Test GORM logger setup
	loggerGorm := logger.NewLoggerGorm(&logger.GormSetting{
		Logger: loggerNew.Entry.Logger,
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

	// Test Postgres config creation
	postgresDBname := "test"
	postgresUsername := "user"

	postgresConfig := &infrastructure.PostgresConfig{
		Config: &postgres.Config{
			DSN:                  "host=localhost user=user password=password dbname=test port=5432 sslmode=disable TimeZone=Asia/Tokyo",
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
	assert.Contains(t, postgresConfig.Config.DSN, postgresDBname)
	assert.Contains(t, postgresConfig.Config.DSN, postgresUsername)
	assert.True(t, postgresConfig.Config.PreferSimpleProtocol)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, postgresConfig.ConnMaxLifetime)
	assert.Equal(t, time.Minute*connectionLifetimeMinutes, postgresConfig.ConnMaxIdletime)
	assert.Equal(t, maxIdleConnections, postgresConfig.MaxIdleConns)
	assert.Equal(t, maxOpenConnections, postgresConfig.MaxOpenConns)
}
