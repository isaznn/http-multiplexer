.PHONY:run
run:
	SRV_HOST=0.0.0.0 SRV_PORT=8080 go run ./cmd/muxer/main.go

.PHONY:build
build:
	go build -o muxer ./cmd/muxer/main.go

.PHONY:test
test:
	go test ./internal/handler ./internal/httpreq
