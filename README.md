# 増分ゲームAPI

My first go project is a community incremental game that can be played using a
commandline application. Includes both a Server with JSON RPC, and a CLI for 
interacting.

This application is a little tongue in cheek and is very much a "toy project"
as the goal was to have something simple that exposed me to the top to bottom
of a typical go application stack.

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

Run `make devu` to spin up the containers locally
Run `make devd` to delete the containers

Can test using curl, the cli or the TUI to play.

## Features for Learning Go

- [x] API
    - [x] /count
    - [x] /increment
    - [x] /motd to display a message of the day
    - [x] /register
    - [x] /healthcheck
- [x] Containerized services
- [x] Prometheus integration for metrics collection
- [x] Grafana local service for observability dashboards
- [x] DB tasks
    - [x] implement Schema
    - [x] API Key registration/management
    - [x] sqlc functions for each API request
    - [x] migration service
- [x] Dedicated CLI implementing the API methods
- [x] Improve error handling ergonomics in API and CLI
- [ ] Improve logging (migrate to slog)
- [x] TUI frontend using bubbletea library
- [ ] Configure Github actions to build the custom images and load them onto docker.
- [ ] Make a companion repository to help facilitate deployment on a VPS

## Potential Improvements for Further Learning
- [ ] Rate Limiting
- [ ] Request Logging Beyond Telemetry
- [ ] Flesh out the TUI functionality
