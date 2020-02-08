# Backend

## Prerequisites

* Golang 1.13+

## How to run unit tests

    go test ./...  -count=1

## How to start DB

    docker-compose up -d mysql

## How to run integration tests

1. Start DB as described above
2. `go test ./...  -tags sql -count=1`
