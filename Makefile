build:
	go build app/main.go

build-linux:
 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o firboom_linux app/main.go

run:
	go build app/main.go && ./main

swagger:
	swag init -g app/main.go -o app/docs

check:
	go vet ./...