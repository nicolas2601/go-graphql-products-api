## ADDED Requirements

### Requirement: Create product use case
The system SHALL create a product through the application layer, assigning an injected id and creation timestamp, and MUST validate the invariants before persisting.

#### Scenario: Valid creation
- **WHEN** the create use case receives valid data
- **THEN** the product is persisted with the injected id and createdAt and returned

#### Scenario: Invalid creation
- **WHEN** the create use case receives data that violates an invariant
- **THEN** the corresponding domain error is returned and nothing is persisted

### Requirement: Query and mutate products
The system SHALL provide get, list, update and delete use cases over the product repository, propagating `ErrProductNotFound` for unknown ids.

#### Scenario: Update existing product
- **WHEN** the update use case receives a known id and a new name or price
- **THEN** the provided fields are applied, revalidated and persisted

#### Scenario: Operate on unknown id
- **WHEN** get, update or delete is called with an unknown id
- **THEN** the system returns `ErrProductNotFound`
