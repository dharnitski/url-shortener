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

## Run App

1. Build Client app as described in `/client/README.md`
2. Start DB as described above
3. `DB_CONNECTION_STRING="root@tcp(localhost:3306)/url-shortener" UI_DOMAIN="http://localhost:8080/"  go run main.go`
