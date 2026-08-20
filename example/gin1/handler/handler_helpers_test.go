//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var errTestCreateTable = errors.New("create table failed")

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestLogger() *logger.Logger {
	return logger.NewLogger(logrus.New())
}

func invokeHandler(t *testing.T, fn func(*gin.Context)) (int, string) {
	t.Helper()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	fn(c)
	return rec.Code, rec.Body.String()
}

func newMySQLGormDB(t *testing.T) *gorm.DB {
	t.Helper()

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock.New() error: %v", err)
	}
	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec("CREATE TABLE `user` (`name` longtext,`email` longtext,`id` bigint AUTO_INCREMENT,PRIMARY KEY (`id`))").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `user` (`name`,`email`) VALUES (?,?)").
		WithArgs("test", "test@test.com").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT * FROM `user` LIMIT ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "test", "test@test.com"))
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 0;").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS `user` CASCADE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS = 1;").WillReturnResult(sqlmock.NewResult(0, 0))

	db, err := gorm.Open(
		mysql.New(mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true}),
		infrastructure.GetDefaultGormConfig(),
	)
	if err != nil {
		t.Fatalf("gorm.Open(mysql) error: %v", err)
	}
	return db
}

func newMySQLGormDBCreateTableError(t *testing.T) *gorm.DB {
	t.Helper()

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock.New() error: %v", err)
	}
	mock.ExpectExec("CREATE TABLE `user` (`name` longtext,`email` longtext,`id` bigint AUTO_INCREMENT,PRIMARY KEY (`id`))").
		WillReturnError(errTestCreateTable)

	db, err := gorm.Open(
		mysql.New(mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true}),
		infrastructure.GetDefaultGormConfig(),
	)
	if err != nil {
		t.Fatalf("gorm.Open(mysql) error: %v", err)
	}
	return db
}

func newPostgresGormDB(t *testing.T) *gorm.DB {
	t.Helper()

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock.New() error: %v", err)
	}
	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec("CREATE TABLE \"user\" (\"name\" text,\"email\" text,\"id\" bigserial,PRIMARY KEY (\"id\"))").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"user\" (\"name\",\"email\") VALUES ($1,$2)").
		WithArgs("test", "test@test.com").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT * FROM \"user\" LIMIT $1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "test", "test@test.com"))
	mock.ExpectExec("DROP TABLE IF EXISTS \"user\" CASCADE").WillReturnResult(sqlmock.NewResult(0, 0))

	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: mockDB, PreferSimpleProtocol: true}),
		infrastructure.GetDefaultGormConfig(),
	)
	if err != nil {
		t.Fatalf("gorm.Open(postgres) error: %v", err)
	}
	return db
}

func newPostgresGormDBCreateTableError(t *testing.T) *gorm.DB {
	t.Helper()

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock.New() error: %v", err)
	}
	mock.ExpectExec("CREATE TABLE \"user\" (\"name\" text,\"email\" text,\"id\" bigserial,PRIMARY KEY (\"id\"))").
		WillReturnError(errTestCreateTable)

	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: mockDB, PreferSimpleProtocol: true}),
		infrastructure.GetDefaultGormConfig(),
	)
	if err != nil {
		t.Fatalf("gorm.Open(postgres) error: %v", err)
	}
	return db
}
