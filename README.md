# ![RealWorld Example App](project-logo.png)

![CI](https://github.com/neo4j-examples/golang-neo4j-realworld-example/workflows/Go/badge.svg)

> ### Neo4j & Golang (using the Neo4j driver) codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

## Prerequisites

Make sure to install a recent [Golang](https://golang.org/) version.

## Build

As simple as:

```
go build ./cmd/conduit
```

## Run

make sure to configure the application to target your specific Neo4j instance.
All settings are mandatory.

| Environment variable  | Description |
| --------------------- | ----------- |
| NEO4J_URI             | [Connection URI](https://neo4j.com/docs/driver-manual/current/client-applications/#driver-connection-uris) of the instance (e.g. `bolt://localhost`, `neo4j+s://example.org`) |
| NEO4J_USERNAME        | Username of the account to connect with (must have read & write permissions) |
| NEO4J_PASSWORD        | Password of the account to connect with (must have read & write permissions)|

> configure in .zshrc 

```
go run ./cmd/conduit/
```

And exercise the application with [Postman](https://www.postman.com/) (see [the collection file](./Conduit.postman_collection.json)) at `localhost:3000`.