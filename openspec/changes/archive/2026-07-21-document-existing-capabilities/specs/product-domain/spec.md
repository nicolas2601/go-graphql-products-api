## ADDED Requirements

### Requirement: Product invariants
The system SHALL construct a product only when its business invariants hold: the name MUST NOT be empty, the price MUST be greater than zero and the stock MUST NOT be negative.

#### Scenario: Valid product
- **WHEN** a product is built with a non-empty name, a price greater than zero and a non-negative stock
- **THEN** the product is created without error

#### Scenario: Empty name
- **WHEN** a product is built with an empty or blank name
- **THEN** the system returns `ErrInvalidName`

#### Scenario: Non-positive price
- **WHEN** a product is built with a price less than or equal to zero
- **THEN** the system returns `ErrInvalidPrice`

#### Scenario: Negative stock
- **WHEN** a product is built with a negative stock
- **THEN** the system returns `ErrInvalidStock`

### Requirement: Typed domain errors
The system SHALL expose typed sentinel errors for domain failures so that outer layers can detect them with `errors.Is`.

#### Scenario: Not found error
- **WHEN** an operation cannot find a product by id
- **THEN** the system returns `ErrProductNotFound`
