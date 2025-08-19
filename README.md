# Podcast API

## .env file

```
PORT=
DB_URL=postgres://postgres:username@domain:port/rssagg?sslmode=disable
KEY=
ALLOWED_ORIGINS=http://*,https://*
SCRAPE_FREQUENCY=1
SCRAPE_CONCURRENCY=10
```

## Database Migrations

* [Goose](https://github.com/pressly/goose)

```
export GOOSE_MIGRATION_DIR=./sql/schema
export GOOSE_DRIVE=postgres
export GOOSE_DBSTRING=postgres://postgres:username@domain:port/rssagg?sslmode=disable
goose postgres up
```

## SQL Generator

* [sqlc](https://docs.sqlc.dev/en/latest/)

```
sqlc compile
sqlc generate
```

## Build

```
go mod tidy
go build
./rssagg
```

