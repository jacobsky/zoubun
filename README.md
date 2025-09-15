# 増分ゲームAPI

My first go project is a community incremental game that can be played using a
commandline application. Includes both a Server with JSON RPC, Hypermedia App,
and CLI (for power gamers).

This application is a little tongue in cheek, but should demonstrate my general
understanding and learning of go.

## How to run

Define your .env file as per the following sample

```sh
PORT=3000
LOG_LEVEL=TRACE
POSTGRES_HOST="zoubun-postgres"
POSTGRES_DB="zoubun-db"
POSTGRES_PASSWORD="hirake5ma"
POSTGRES_USER="user"
```

Run `docker-compose up`

## Features In Development

- [x] API
- [x] Containerized services
- [ ] Hypermedia App
- [ ] DB backed counting and registration
- [ ] User Registration - Consider Discord OAuth for verification.
- [ ] Contribution tracking page
- [ ] Dedicated CLI
