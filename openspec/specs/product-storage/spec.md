# product-storage Specification

## Purpose
Define el contrato de persistencia de productos y su implementacion en memoria segura para uso concurrente.
## Requirements
### Requirement: Product repository contract
The system SHALL persist products through a repository interface defined in the domain layer, returning `ErrProductNotFound` for unknown ids and `ErrProductAlreadyExists` when creating a duplicate id.

#### Scenario: Persist and retrieve
- **WHEN** a product is created and then requested by id
- **THEN** the stored product is returned

#### Scenario: Duplicate id
- **WHEN** a product is created with an id that already exists
- **THEN** the system returns `ErrProductAlreadyExists`

### Requirement: Concurrency-safe in-memory storage
The in-memory repository SHALL be safe for concurrent use and MUST return the product list in a deterministic order.

#### Scenario: Concurrent access
- **WHEN** multiple goroutines read and write the repository concurrently
- **THEN** no data race occurs

#### Scenario: Deterministic listing
- **WHEN** the products are listed
- **THEN** they are returned ordered by id

### Requirement: PostgreSQL persistence
The system SHALL provide a PostgreSQL-backed implementation of the product repository, selectable via `REPO_DRIVER=postgres`, applying its schema through a migration.

#### Scenario: Persist and retrieve
- **WHEN** the PostgreSQL repository creates a product and then requests it by id
- **THEN** the stored product is returned from the database

#### Scenario: Duplicate id
- **WHEN** the PostgreSQL repository creates a product with an id that already exists
- **THEN** it returns `ErrProductAlreadyExists`

#### Scenario: Unknown id
- **WHEN** the PostgreSQL repository is asked for, updates or deletes an unknown id
- **THEN** it returns `ErrProductNotFound`

