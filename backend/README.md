## PackServer

API server for ParkEasy

### Getting started

Start a local cluster for development

    docker compose -f compose.yaml up --build

The API server is exposed on port `8080`.

Try a simple API request

    curl http://localhost:8080/greeting/world

#### Hot code reloading

Hot code reloading is provided by [Air](https://github.com/air-verse/air) and can be run via

    docker compose -f compose.air.yaml up --build

### Architecture

TBD
