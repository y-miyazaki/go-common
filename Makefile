build-initial:
	GOOS=darwin GOARCH=amd64 go build example/initial/main.go
build-gin1:
	GOOS=darwin GOARCH=amd64 go build example/gin1/main.go
build-gin2:
	GOOS=darwin GOARCH=amd64 go build example/gin2/main.go
