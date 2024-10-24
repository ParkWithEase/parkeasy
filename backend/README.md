## PackServer

API server for ParkEasy

### Getting started

Start a local cluster for development

    cp example.env .env
    docker compose -f compose.yaml up --build

The API server is exposed on port `8080`.

The API documentation server can then be reached at `http://localhost:8080/docs`

#### Hot code reloading

Hot code reloading is provided by [Air] and can be run via

    docker compose -f compose.yaml -f compose.air.yaml up --build

[Air] would then react automatically upon changes in the source code and reloads the server automatically, allowing for a quick feedback loop.

[Air]: https://github.com/air-verse/air

### DB model

This server uses [Bob](https://bob.stephenafamo.com/) for building and running SQL queries.

In order to get the most out of the library, a model must be generated from the database. This has to be done whenever the database schemas are modified. To do this, run:

    go run ./tools/modelgen

### Tests

Tests can be run via

    go test ./...

By default, this would only run unit tests. To run integration tests as well, run:

    INTEGRATION=1 go test ./...

To run only integration tests, run:

    INTEGRATION=1 go test -run Integration ./...

For integration tests to run successfully, Docker should be installed and running.

### Architecture

TBD
