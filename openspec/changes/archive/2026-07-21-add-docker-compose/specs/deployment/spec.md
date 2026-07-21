## ADDED Requirements

### Requirement: Containerized deployment
The system SHALL be runnable as a container image and orchestrated via docker-compose together with a PostgreSQL service.

#### Scenario: Compose brings up the API with PostgreSQL
- **WHEN** the operator runs `docker compose up`
- **THEN** the API container starts after the PostgreSQL service is healthy and serves GraphQL persisting to that database

#### Scenario: Container health check
- **WHEN** the API container is running
- **THEN** its health check reports the container as healthy through the `/healthz` endpoint

### Requirement: Non-root container
The container image SHALL run the application as a non-root user.

#### Scenario: Runtime user
- **WHEN** the image is built and run
- **THEN** the application process runs as a non-root user
