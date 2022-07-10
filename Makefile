.PHONY:run
run:
	SRV_PORT=8080 go run ./cmd/muxer/main.go

.PHONY:build
build:
	go build -o muxer ./cmd/muxer/main.go
