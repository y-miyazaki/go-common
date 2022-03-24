package utils

import "strings"

// GetMySQLDsn creates a MySQL DSN string.
// about DSN document.
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func GetMySQLDsn(username, password, server, port, dbname, parameter string) string {
	if !strings.HasPrefix(parameter, "?") {
		parameter = "?" + parameter
	}
	return username + ":" + password + "@tcp(" + server + ":" + port + ")/" + dbname + parameter
}
