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

Run `make devu` to spin up the containers locally
Run `make devd` to delete the containers


## Features In Development

- [ ] API
    - [x] /count
    - [x] /increment
    - [x] /motd to display a message of the day
    - [ ] /register
    - [ ] /verification
    - [ ] /healthcheck
    - [ ] etc...
- [x] Containerized services
- [x] Prometheus integration for metrics collection
- [x] Grafana local service for observability dashboards
- [ ] DB tasks
    - [ ] implement Schema
    - [ ] API Key registration/management
    - [ ] sqlc functions for each API request
    - [x] migration service
- [ ] Dedicated CLI implementing the API methods
- [ ] Hypermedia based counting visualization
- [ ] Configure Github actions to build the custom images and load them onto docker.
- [ ] (Stretch goal) Add in simple email verification service to verify newly registered accounts
- [ ] Make a companion repository to help facilitate deployment on a VPS
