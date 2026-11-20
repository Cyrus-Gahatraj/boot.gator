# gator

An RSS feed aggregator CLI.

## Requirements

- Go
- Postgres

## Install

```sh
go install github.com/Cyrus-Gahatraj/boot.gator@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`. The installed binary is named `boot.gator` (`alias gator=boot.gator` if you prefer `gator`).

## Setup

Create a Postgres database, then create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Apply the schema (requires [goose](https://github.com/pressly/goose)):

```sh
git clone https://github.com/Cyrus-Gahatraj/boot.gator.git
cd boot.gator
GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://user:password@localhost:5432/gator?sslmode=disable" goose -dir sql/schema up
```

## Usage

```sh
boot.gator register alice          # create a user and log in
boot.gator login alice             # switch users
boot.gator addfeed "My Feed" https://example.com/rss
boot.gator follow https://example.com/rss
boot.gator agg 1m                  # scrape followed feeds every minute
boot.gator browse 5                # show the 5 newest posts (default: 2)
boot.gator feeds                   # list all feeds
boot.gator following               # list feeds you follow
boot.gator users                   # list users
```
