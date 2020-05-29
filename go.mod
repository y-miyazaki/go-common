module github.com/y-miyazaki/go-common

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/y-miyazaki/go-common/pkg/errors => ./pkg/errors
