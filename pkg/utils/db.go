package utils

// GetMySqlDsn creates a MySQL DSN string.
// about DSN document.
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
func GetMySqlDsn(username, password, server, port, dbname, parameter string) string {
	if !strings.HasPrefix(parameter, "?"){
		parameter = "?" + parameter
	}
	charset=utf8mb4&parseTime=True&loc=Local
	return username + ":" + password + "@tcp(" + server + ":" + port + ")/" + dbname + parameter
}
