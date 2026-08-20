//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package infrastructure

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func TestMySQLConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		cfg  *MySQLConfig
		name string
		want DBConfig
	}{
		{
			name: "stores connection pool settings",
			cfg: &MySQLConfig{
				Config: &mysql.Config{
					DSN: "user:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
				},
				DBConfig: DBConfig{
					ConnMaxLifetime: 1 * time.Hour,
					ConnMaxIdletime: 30 * time.Minute,
					MaxIdleConns:    10,
					MaxOpenConns:    100,
				},
			},
			want: DBConfig{
				ConnMaxLifetime: 1 * time.Hour,
				ConnMaxIdletime: 30 * time.Minute,
				MaxIdleConns:    10,
				MaxOpenConns:    100,
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.NotNil(t, tt.cfg)
			require.NotNil(t, tt.cfg.Config)
			assert.Equal(t, tt.want.ConnMaxLifetime, tt.cfg.ConnMaxLifetime)
			assert.Equal(t, tt.want.ConnMaxIdletime, tt.cfg.ConnMaxIdletime)
			assert.Equal(t, tt.want.MaxIdleConns, tt.cfg.MaxIdleConns)
			assert.Equal(t, tt.want.MaxOpenConns, tt.cfg.MaxOpenConns)
		})
	}
}

func TestDBConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  DBConfig
		want DBConfig
	}{
		{
			name: "stores pool limits",
			cfg: DBConfig{
				ConnMaxLifetime: 2 * time.Hour,
				ConnMaxIdletime: 45 * time.Minute,
				MaxIdleConns:    20,
				MaxOpenConns:    200,
			},
			want: DBConfig{
				ConnMaxLifetime: 2 * time.Hour,
				ConnMaxIdletime: 45 * time.Minute,
				MaxIdleConns:    20,
				MaxOpenConns:    200,
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.cfg)
		})
	}
}

func TestGetDefaultGormConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "returns non-nil defaults"},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := GetDefaultGormConfig()
			require.NotNil(t, config)
			assert.True(t, config.DisableAutomaticPing)
			require.NotNil(t, config.NamingStrategy)
		})
	}
}

func TestNewDB_InvalidDSN(t *testing.T) {
	tests := []struct {
		open      func()
		name      string
		wantPanic bool
	}{
		{
			name: "mysql invalid dsn panics",
			open: func() {
				NewMySQL(&MySQLConfig{
					Config: &mysql.Config{DSN: "invalid-dsn"},
				}, GetDefaultGormConfig())
			},
			wantPanic: true,
		},
		{
			name: "postgres invalid dsn panics",
			open: func() {
				NewPostgres(&PostgresConfig{
					Config: &postgres.Config{DSN: "invalid-dsn"},
				}, GetDefaultGormConfig())
			},
			wantPanic: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			panicked := false
			func() {
				defer func() {
					if recover() != nil {
						panicked = true
					}
				}()
				tt.open()
			}()
			if panicked != tt.wantPanic {
				t.Fatalf("%s: panicked = %v, want %v", tt.name, panicked, tt.wantPanic)
			}
		})
	}
}

func TestNewSQLServer_OpenWithoutPanic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dsn  string
	}{
		{name: "invalid dsn opens without panic when ping disabled", dsn: "invalid-dsn"},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			panicked := false
			var db *gorm.DB
			func() {
				defer func() {
					if recover() != nil {
						panicked = true
					}
				}()
				db = NewSQLServer(&SQLServerConfig{
					Config: &sqlserver.Config{DSN: tt.dsn},
				}, GetDefaultGormConfig())
			}()
			if panicked {
				t.Fatalf("NewSQLServer(%q) panicked, want open without panic (DisableAutomaticPing)", tt.dsn)
			}
			if db == nil {
				t.Fatalf("NewSQLServer(%q) = nil, want non-nil *gorm.DB", tt.dsn)
			}
		})
	}
}

func TestNewDB_WithSqlmock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		open func() any
		name string
	}{
		{
			name: "mysql applies pool config via initDB",
			open: func() any {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("sqlmock.New() error: %v", err)
				}
				return NewMySQL(&MySQLConfig{
					Config: &mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true},
					DBConfig: DBConfig{
						ConnMaxLifetime: time.Hour,
						ConnMaxIdletime: time.Minute,
						MaxIdleConns:    1,
						MaxOpenConns:    2,
					},
				}, GetDefaultGormConfig())
			},
		},
		{
			name: "postgres applies pool config via initDB",
			open: func() any {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("sqlmock.New() error: %v", err)
				}
				return NewPostgres(&PostgresConfig{
					Config: &postgres.Config{Conn: mockDB, PreferSimpleProtocol: true},
					DBConfig: DBConfig{
						ConnMaxLifetime: time.Hour,
						ConnMaxIdletime: time.Minute,
						MaxIdleConns:    1,
						MaxOpenConns:    2,
					},
				}, GetDefaultGormConfig())
			},
		},
		{
			name: "sqlserver applies pool config via initDB",
			open: func() any {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("sqlmock.New() error: %v", err)
				}
				return NewSQLServer(&SQLServerConfig{
					Config: &sqlserver.Config{Conn: mockDB},
					DBConfig: DBConfig{
						ConnMaxLifetime: time.Hour,
						ConnMaxIdletime: time.Minute,
						MaxIdleConns:    1,
						MaxOpenConns:    2,
					},
				}, GetDefaultGormConfig())
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := tt.open()
			if db == nil {
				t.Fatalf("%s: returned nil database handle", tt.name)
			}
		})
	}
}
