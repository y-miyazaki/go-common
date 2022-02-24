package infrastructure

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// MySQLConfig sets configurations.
type MySQLConfig struct {
	Config *mysql.Config
	DBConfig
}

// PostgresConfig sets configurations.
type PostgresConfig struct {
	Config *postgres.Config
	DBConfig
}

// SQLServerConfig sets configurations.
type SQLServerConfig struct {
	Config *sqlserver.Config
	DBConfig
}

// DBConfig set configurations.
type DBConfig struct {
	// ConnMaxLifetime sets max life time(sec)
	ConnMaxLifetime time.Duration
	// ConnMaxIdletime sets max idle time(sec)
	ConnMaxIdletime time.Duration
	// MaxIdleConns sets idle connection
	MaxIdleConns int
	// MaxOpenConns sets max connection
	MaxOpenConns int
}

// NewMySQL returns an gorm db instance.
func NewMySQL(c *MySQLConfig, gc *gorm.Config) *gorm.DB {
	i := mysql.New(*c.Config)
	db, err := gorm.Open(i, gc)
	if err != nil {
		panic(fmt.Sprintf("can't open error. %v", err))
	}
	initDB(db, c.DBConfig)
	return db
}

// NewPostgres returns an gorm db instance.
func NewPostgres(c *PostgresConfig, gc *gorm.Config) *gorm.DB {
	i := postgres.New(*c.Config)
	db, err := gorm.Open(i, gc)
	if err != nil {
		panic(fmt.Sprintf("can't open error. %v", err))
	}
	initDB(db, c.DBConfig)
	return db
}

// NewSQLServer returns an gorm db instance.
func NewSQLServer(c *SQLServerConfig, gc *gorm.Config) *gorm.DB {
	i := sqlserver.New(*c.Config)
	db, err := gorm.Open(i, gc)
	if err != nil {
		panic(fmt.Sprintf("can't open error. %v", err))
	}
	initDB(db, c.DBConfig)
	return db
}

// GetDefaultGormConfig get default config.
func GetDefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func initDB(db *gorm.DB, dbConfig DBConfig) {
	d, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("can't get db instance error. %v", err))
	}
	d.SetConnMaxIdleTime(dbConfig.ConnMaxIdletime * time.Second)
	d.SetConnMaxLifetime(dbConfig.ConnMaxLifetime * time.Second)
	d.SetMaxIdleConns(dbConfig.MaxIdleConns)
	d.SetMaxOpenConns(dbConfig.MaxOpenConns)
}
