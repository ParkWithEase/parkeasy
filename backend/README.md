## PackServer

API server for ParkEasy

### Getting started

Start a local cluster for development

    cp example.env .env
    docker compose -f compose.yaml up --build

The API server is exposed on port `8080`.

The documentation server can then be reached at `http://localhost:8080/docs`

#### Hot code reloading

Hot code reloading is provided by [Air](https://github.com/air-verse/air) and can be run via

    docker compose -f compose.air.yaml up --build

### Tests

Tests can be run via

    go test ./...

Currently no distinction exists between unit and integration tests

### Architecture

TBD
