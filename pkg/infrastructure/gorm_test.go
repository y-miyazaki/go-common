package infrastructure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
)

func TestNewMySQL(t *testing.T) {
	// This test would require a real MySQL database connection
	// For now, just test the config struct creation
	config := &MySQLConfig{
		Config: &mysql.Config{
			DSN: "user:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		},
		DBConfig: DBConfig{
			ConnMaxLifetime: 1 * time.Hour,
			ConnMaxIdletime: 30 * time.Minute,
			MaxIdleConns:    10,
			MaxOpenConns:    100,
		},
	}

	assert.NotNil(t, config)
	assert.NotNil(t, config.Config)
	assert.Equal(t, 1*time.Hour, config.ConnMaxLifetime)
	assert.Equal(t, 30*time.Minute, config.ConnMaxIdletime)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 100, config.MaxOpenConns)
}

func TestDBConfig(t *testing.T) {
	config := DBConfig{
		ConnMaxLifetime: 2 * time.Hour,
		ConnMaxIdletime: 45 * time.Minute,
		MaxIdleConns:    20,
		MaxOpenConns:    200,
	}

	assert.Equal(t, 2*time.Hour, config.ConnMaxLifetime)
	assert.Equal(t, 45*time.Minute, config.ConnMaxIdletime)
	assert.Equal(t, 20, config.MaxIdleConns)
	assert.Equal(t, 200, config.MaxOpenConns)
}
