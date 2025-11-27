# 🦥 Sloth Monitor

This project is a small **lazy** tool to monitor Source Query ports of servers in a database and update their online player counts. It can be used as a long-running backend service for some server list frontends.

## Usage

1. Create `.env` file with database details.

```sh
DB_USER=username
DB_PASS=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=database
DB_TABLE=servers
POLL_SECONDS=30
```

2. Use `go run .` to run the app.