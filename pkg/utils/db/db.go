// Package db provides database utility functions for connection string generation.
package db

import (
	"fmt"
	"strings"
)

// GetMySQLDsn creates a MySQL DSN string.
// about DSN document.
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func GetMySQLDsn(username, password, server, port, db, parameter string) string {
	finalParameter := parameter
	if !strings.HasPrefix(parameter, "?") {
		finalParameter = "?" + parameter
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", username, password, server, port, db, finalParameter)
}

// GetPostgresDsn creates a Postgres DSN string.
func GetPostgresDsn(user, password, host, port, dbname, parameter string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s", user, password, host, port, dbname, parameter)
}
