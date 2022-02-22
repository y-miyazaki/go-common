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

// MySQLConfigSetting sets configurations.
type MySQLConfigSetting struct {
	Config *mysql.Config
	DBConfig
}

// PostgresConfigSetting sets configurations.
type PostgresConfigSetting struct {
	Config *postgres.Config
	DBConfig
}

// SQLServerConfigSetting sets configurations.
type SQLServerConfigSetting struct {
	Config *sqlserver.Config
	DBConfig
}

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
func NewMySQL(c *MySQLConfigSetting, gc *gorm.Config) *gorm.DB {
	i := mysql.New(*c.Config)
	db, err := gorm.Open(i, gc)
	if err != nil {
		panic(fmt.Sprintf("can't open error. %v", err))
	}
	initDB(db, c.DBConfig)
	return db
}

// NewPostgres returns an gorm db instance.
func NewPostgres(c *PostgresConfigSetting, gc *gorm.Config) *gorm.DB {
	i := postgres.New(*c.Config)
	db, err := gorm.Open(i, gc)
	if err != nil {
		panic(fmt.Sprintf("can't open error. %v", err))
	}
	initDB(db, c.DBConfig)
	return db
}

// NewSQLServer returns an gorm db instance.
func NewSQLServer(c *SQLServerConfigSetting, gc *gorm.Config) *gorm.DB {
	i := sqlserver.New(*c.Config)
	db, err := gorm.Open(i, gc)
	if err != nil {
		panic(fmt.Sprintf("can't open error. %v", err))
	}
	initDB(db, c.DBConfig)
	return db
}

func GetDefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func initDB(db *gorm.DB, dbConfig DBConfig) {
	sqldb, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("can't get db instance error. %v", err))
	}
	sqldb.SetConnMaxIdleTime(dbConfig.ConnMaxIdletime * time.Second)
	sqldb.SetConnMaxLifetime(dbConfig.ConnMaxLifetime * time.Second)
	sqldb.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqldb.SetMaxOpenConns(dbConfig.MaxOpenConns)
}
