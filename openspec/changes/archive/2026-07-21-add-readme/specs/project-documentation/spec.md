## ADDED Requirements

### Requirement: Project README
The project SHALL provide a README that documents setup, configuration, both run modes (in-memory and PostgreSQL), the GraphQL API, and the design decisions.

#### Scenario: Running from the README
- **WHEN** a developer follows the README instructions
- **THEN** they can run the API both in in-memory mode and with PostgreSQL via docker-compose

#### Scenario: Configuration reference
- **WHEN** a developer needs to configure the application
- **THEN** the README lists every environment variable with its default and purpose
