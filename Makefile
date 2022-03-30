build-mysql:
	GOOS=darwin GOARCH=amd64 go build example/mysql/main.go
build-postgres:
	GOOS=darwin GOARCH=amd64 go build example/postgres/main.go
build-s3:
	GOOS=darwin GOARCH=amd64 go build example/s3/main.go
build-gin1:
	GOOS=darwin GOARCH=amd64 go build example/gin1/main.go
build-gin2:
	GOOS=darwin GOARCH=amd64 go build example/gin2/main.go
