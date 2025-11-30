format:
	gofmt -w .

build:
	go build -o main .

air_start:
	~/go/bin/air

start:
	go run main.go