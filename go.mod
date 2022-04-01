module github.com/y-miyazaki/go-common

go 1.16

replace github.com/y-miyazaki/go-common v0.0.0 => ./

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-sdk-go v1.43.3
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/mock v1.6.0
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/rivo/uniseg v0.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/slack-go/slack v0.10.2
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/exp v0.0.0-20220218215828-6cf2b201936e
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gorm.io/driver/mysql v1.3.2
	gorm.io/driver/postgres v1.3.1
	gorm.io/driver/sqlserver v1.3.1
	gorm.io/gorm v1.23.1
	gorm.io/plugin/soft_delete v1.1.0
)
