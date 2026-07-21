## ADDED Requirements

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
