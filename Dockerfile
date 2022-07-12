FROM golang:1.18.3-alpine
ENV CGO_ENABLED=0

WORKDIR /app
COPY . /app

RUN go build -o muxer ./cmd/muxer/main.go

CMD [ "./muxer" ]
