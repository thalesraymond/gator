# Gator

Gator is a guided Boot.dev study project: a CLI feed aggregator written in Go.

It lets you:

- manage local users,
- add and follow RSS feeds,
- periodically fetch feed entries into Postgres,
- browse the latest posts from followed feeds.

## Tech Stack

- Go 1.25+
- PostgreSQL
- [goose](https://github.com/pressly/goose) for migrations
- [sqlc](https://sqlc.dev/) for type-safe DB access

## Features

- User management: register, login, list users, reset users
- Feed management: add feeds and list all feeds
- Follow system: follow/unfollow feeds and list followed feeds
- Aggregation worker: polls feeds and stores posts
- Browsing: read latest posts from followed feeds

## Project Structure

- `cmd/`: CLI command handlers
- `internal/config/`: reads and writes `~/.gatorconfig.json`
- `internal/database/`: sqlc-generated query layer
- `internal/feed/`: RSS fetch and parse helpers
- `sql/schema/`: goose migrations
- `sql/queries/`: sqlc query definitions

## Prerequisites

Install these before running the app:

1. Go
2. PostgreSQL
3. goose CLI

Optional (only needed if you plan to regenerate DB code):

4. sqlc CLI

## Setup

### 1) Create the config file

Create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

### 2) Create local migration env

Create a `.env` file in the project root for goose:

```bash
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/gator?sslmode=disable
GOOSE_MIGRATION_DIR=./sql/schema
```

### 3) Run migrations

Use the helper script:

```bash
./goose.sh up
```

Or run goose directly:

```bash
export $(cat .env | xargs) && goose up
```

### 4) Build and run

```bash
go build -o gator
./gator
```

You can also run commands directly with:

```bash
go run . <command> [args...]
```

## Commands

### User commands

- `register <username>`
- `login <username>`
- `users`
- `reset`

### Feed commands

- `addfeed <name> <url>` (requires logged-in user)
- `feeds` (requires logged-in user)

### Follow commands

- `follow <url>` (requires logged-in user)
- `following` (requires logged-in user)
- `unfollow <url>` (requires logged-in user)

### Aggregation and browsing

- `agg` runs the feed fetch worker loop
- `browse [limit]` (requires logged-in user, default limit is 2)

## Quick Start Flow

After setup and migrations:

```bash
go run . register thales
go run . addfeed "Hacker News" "https://hnrss.org/frontpage"
go run . agg
```

In another terminal:

```bash
go run . browse 10
```

## Notes

- `agg` is intentionally long-running.
- Feed URLs and post URLs are unique in the database.
- `addfeed` also auto-follows the newly created feed for the current user.

## Study Context

This project was built as part of the Boot.dev curriculum and is intended as a learning-focused implementation of a Go + Postgres CLI app.
