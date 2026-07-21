# server-runtime Specification

## Purpose
TBD - created by archiving change add-server-wiring. Update Purpose after archive.
## Requirements
### Requirement: HTTP server startup
The system SHALL provide an executable that starts an HTTP server on the configured port, exposing the GraphQL API and a health check.

#### Scenario: Health check
- **WHEN** a client requests `GET /healthz`
- **THEN** the server responds with HTTP 200 and a JSON status body

#### Scenario: GraphQL endpoint
- **WHEN** a client sends a GraphQL operation to `POST /query`
- **THEN** the server executes it against the resolvers and returns the result

### Requirement: Repository selection by configuration
The system SHALL select the product repository from the `REPO_DRIVER` configuration, defaulting to the in-memory repository.

#### Scenario: In-memory driver
- **WHEN** `REPO_DRIVER` is `memory`
- **THEN** the system wires the in-memory repository

#### Scenario: Unsupported driver
- **WHEN** `REPO_DRIVER` is a value that is not yet supported
- **THEN** the system fails to start with an explicit error

### Requirement: Graceful shutdown
The system SHALL shut down gracefully when it receives an interrupt or termination signal, letting in-flight requests finish within a timeout.

#### Scenario: Termination signal
- **WHEN** the process receives SIGINT or SIGTERM
- **THEN** the server stops accepting new connections and shuts down within the configured timeout

