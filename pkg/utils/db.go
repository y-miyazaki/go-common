package utils

import (
	"fmt"
	"strings"
)

// GetMySQLDsn creates a MySQL DSN string.
// about DSN document.
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func GetMySQLDsn(username, password, server, port, db, parameter string) string {
	if !strings.HasPrefix(parameter, "?") {
		parameter = "?" + parameter
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", username, password, server, port, db, parameter)
}

// GetPostgresDsn creates a Postgres DSN string.
// about DSN document.
// https://github.com/go-sql-driver/postgres#dsn-data-source-name
func GetPostgresDsn(user, password, host, port, dbname, parameter string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s", user, password, host, port, dbname, parameter)
}
