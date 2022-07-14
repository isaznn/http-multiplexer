#!/bin/bash

docker build . -t muxer
docker run muxer sh -c "go test -count=10 ./internal/external ./internal/service ./internal/handler"
