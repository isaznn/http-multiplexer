# http-multiplexer

Concurrent multiplexer for HTTP requests using Go 1.18.

Limits (configured via constants in main.go):
- 100 concurrent incoming HTTP requests
- 4 simultaneous outgoing requests per incoming request
- 1 second HTTP client timeout
- 20 urls per request

## Data structure

Request:
```
{
  "urls": ["<url_1>", "<url_2>", ...]
}
```

Response:
```
{
  "error": false,
  "result": {
    "<url_1>": "<Response>",
    "<url_2>": "<Response>",
    ...
  }
}
```

Error:
```
{
  "error": true,
  "errorMessage": "<Error>"
}
```

## Usage

POST /muxer with JSON object containing an array of URLs.

Example:
```
curl -X POST http://localhost:8080/muxer \
    -H 'Content-Type: application/json' \
    -d '{"urls": ["https://jsonplaceholder.typicode.com/posts"]}'
```

## Test & Run

### with make

```
make test
make run
```

### with Docker

```
./docker_test.sh
./docker_run.sh
```

## Build

```
make build
```
Put `muxer` in `./bin` directory
