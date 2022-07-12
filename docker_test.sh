#!/bin/bash

docker build . -t muxer
docker run muxer sh -c "go test -count=1 ./internal/handler ./internal/external"
