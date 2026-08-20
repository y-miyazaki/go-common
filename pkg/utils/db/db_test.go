//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package db

import (
	"testing"
)

func TestGetMySQLDsn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		username  string
		password  string
		server    string
		port      string
		db        string
		parameter string
		want      string
	}{
		{
			name:      "builds tcp dsn with query parameters",
			username:  "test",
			password:  "testpassword",
			server:    "localhost",
			port:      "3306",
			db:        "testdb",
			parameter: "a=1",
			want:      "test:testpassword@tcp(localhost:3306)/testdb?a=1",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetMySQLDsn(tt.username, tt.password, tt.server, tt.port, tt.db, tt.parameter)
			if got != tt.want {
				t.Fatalf("GetMySQLDsn() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetPostgresDsn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		user      string
		password  string
		host      string
		dbname    string
		port      string
		parameter string
		want      string
	}{
		{
			name:      "builds keyword dsn with trailing parameters",
			user:      "test",
			password:  "testpassword",
			host:      "localhost",
			port:      "5432",
			dbname:    "testdb",
			parameter: "a=1",
			want:      "user=test password=testpassword host=localhost port=5432 dbname=testdb a=1",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetPostgresDsn(tt.user, tt.password, tt.host, tt.port, tt.dbname, tt.parameter)
			if got != tt.want {
				t.Fatalf("GetPostgresDsn() = %q, want %q", got, tt.want)
			}
		})
	}
}
