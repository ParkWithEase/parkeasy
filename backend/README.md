## PackServer

API server for ParkEasy

### Getting started

Get a free [geocod.io](https://geocod.io) API key. This is required for listing APIs to work.

Start a local cluster for development

    cp example.env .env
    # edit .env and add your geocod.io API key
    docker compose -f compose.yaml up --build

The API server is exposed on port `8080`.

The API documentation server can then be reached at `http://localhost:8080/docs`

#### Hot code reloading

Hot code reloading is provided by [Air] and can be run via

    docker compose -f compose.yaml -f compose.air.yaml up --build

[Air] would then react automatically upon changes in the source code and reloads the server automatically, allowing for a quick feedback loop.

[Air]: https://github.com/air-verse/air

### Integrated profiler

Runtime profiling data is provided using [`net/http/pprof`](https://pkg.go.dev/net/http/pprof) and can be accessed if `PROVIDER_PORT` is configured. The example environment file will configure this port to 6060.

These endpoints allow for capturing sensitive system data, and control over the garbage collector operations. As such, these endpoints must **not** be exposed to untrusted networks.

For how to make use of this feature, see:

- [`net/http/pprof` Usage examples](https://pkg.go.dev/net/http/pprof#hdr-Usage_examples)
- [Profiling Go Programs](https://go.dev/blog/pprof)

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
