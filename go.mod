module github.com/y-miyazaki/go-common

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-sdk-go v1.45.6
	github.com/aws/aws-sdk-go-v2 v1.18.1
	github.com/aws/aws-sdk-go-v2/config v1.18.26
	github.com/aws/aws-sdk-go-v2/credentials v1.13.25
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.67
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.22.10
	github.com/aws/aws-sdk-go-v2/service/s3 v1.34.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.19.8
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.18.0
	github.com/aws/aws-secretsmanager-caching-go v1.1.0
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/fsnotify/fsnotify v1.6.0
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/pkg/errors v0.9.1
	github.com/rivo/uniseg v0.4.4
	github.com/sirupsen/logrus v1.9.3
	github.com/slack-go/slack v0.12.2
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/driver/sqlserver v1.5.1
	gorm.io/gorm v1.25.1
	gorm.io/plugin/soft_delete v1.2.1
)
