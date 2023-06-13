#   $GOOS     $GOARCH
#   darwin    386
#   darwin    amd64
#   freebsd   386
#   freebsd   amd64
#   freebsd   arm
#   linux     386
#   linux     amd64
#   linux     arm
#   netbsd    386
#   netbsd    amd64
#   netbsd    arm
#   openbsd   386
#   openbsd   amd64
#   plan9     386
#   plan9     amd64
#   windows   386
#   windows   amd64
#   nacl      amd64
#   nacl      386

build-mysql:
	GOOS=linux GOARCH=amd64 go build example/mysql/main.go
build-postgres:
	GOOS=linux GOARCH=amd64 go build example/postgres/main.go
build-s3:
	GOOS=linux GOARCH=amd64 go build example/s3/main.go
build-s3-v2:
	GOOS=linux GOARCH=amd64 go build example/s3_v2/main.go
build-s3-v2-darwin:
	GOOS=linux GOARCH=amd64 go build example/s3_v2/main.go
build-gin1:
	GOOS=linux GOARCH=amd64 go build example/gin1/main.go
build-gin2:
	GOOS=linux GOARCH=amd64 go build example/gin2/main.go
