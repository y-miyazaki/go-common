package db

import (
	"testing"
)

func TestGetMySQLDsn(t *testing.T) {
	type args struct {
		username  string
		password  string
		server    string
		port      string
		db        string
		parameter string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				username:  "test",
				password:  "testpassword",
				server:    "localhost",
				port:      "3306",
				db:        "testdb",
				parameter: "a=1",
			},
			want: "test:testpassword@tcp(localhost:3306)/testdb?a=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMySQLDsn(tt.args.username, tt.args.password, tt.args.server, tt.args.port, tt.args.db, tt.args.parameter); got != tt.want {
				t.Errorf("GetMySQLDsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPostgresDsn(t *testing.T) {
	type args struct {
		user      string
		password  string
		host      string
		dbname    string
		port      string
		parameter string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				user:      "test",
				password:  "testpassword",
				host:      "localhost",
				port:      "5432",
				dbname:    "testdb",
				parameter: "a=1",
			},
			want: "user=test password=testpassword host=localhost port=5432 dbname=testdb a=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPostgresDsn(tt.args.user, tt.args.password, tt.args.host, tt.args.port, tt.args.dbname, tt.args.parameter); got != tt.want {
				t.Errorf("GetPostgresDsn() = %v, want %v", got, tt.want)
			}
		})
	}
}
