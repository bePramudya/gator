# Gator

Gator is a multi-user command-line RSS aggregator written in Go. Users can add and follow feeds, periodically collect posts, and browse the stored articles.

## Requirements

- [Go](https://go.dev/dl/) 1.26.4 or later
- [PostgreSQL](https://www.postgresql.org/download/)
- [Goose](https://github.com/pressly/goose) for database migrations
- Make

## Installation
You can then install gator with:

```sh
go install github.com/bePramudya/gator
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Make sure your Go binary directory is on your `PATH` so the `gator` and `goose` commands are available.

## Config

Create `~/.gatorconfig.json` and set `db_url` to your PostgreSQL connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Set `DB_URL` in the `Makefile` to the same connection string, then apply the database migrations from the repository root:

```makefile
DB_URL := "postgres://user:password@localhost:5432/gator"

migrate-up:
	goose -dir sql/schema postgres $(DB_URL) up

migrate-down:
	goose -dir sql/schema postgres $(DB_URL) down

psql-gator:
	psql $(DB_URL)
```

```sh
make migrate-up
```

## Quick start

Registering a user also makes that user active:

```sh
gator register alice
```

Add a feed. The new feed is automatically followed by its creator:

```sh
gator addfeed "Hacker News" "https://hnrss.org/frontpage"
```

Start collecting posts every minute. The aggregator runs until you stop it with `Ctrl+C`:

```sh
gator agg 1m
```

In another terminal, browse the latest posts:

```sh
gator browse 10
```

## Commands

| Command | Description |
| --- | --- |
| `gator register <name>` | Create a user and make it the active user |
| `gator login <name>` | Switch to an existing user |
| `gator users` | List users and show the active user |
| `gator addfeed <name> <url>` | Add and follow a feed as the active user |
| `gator feeds` | List all feeds, their URLs, and their creators |
| `gator follow <url>` | Follow an existing feed as the active user |
| `gator following` | List feeds followed by the active user |
| `gator unfollow <url>` | Stop following a feed as the active user |
| `gator agg <duration>` | Continuously fetch feeds at intervals such as `30s`, `1m`, or `1h` |
| `gator browse [limit]` | Show recent posts; the default limit is `2` |
| `gator reset` | Delete all users and associated feed data |

> [!WARNING]
> `gator reset` is destructive. Because related records use cascading deletes, it also removes feeds, follows, and posts.

## Development

SQL migrations are in `sql/schema`, SQL queries are in `sql/queries`, and generated database code is in `internal/database`.

After changing a query or schema, regenerate the Go database code with [sqlc](https://sqlc.dev/):

```sh
sqlc generate
```

Check that the project builds and tests pass with:

```sh
go test ./...
```
