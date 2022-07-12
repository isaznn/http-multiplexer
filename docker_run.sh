#!/bin/bash

docker build . -t muxer
docker run -e SRV_HOST=0.0.0.0 -e SRV_PORT=8080 -p 8080:8080 muxer
