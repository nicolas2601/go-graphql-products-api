## ADDED Requirements

### Requirement: List products
The system SHALL return all products through the `products` query.

#### Scenario: Products exist
- **WHEN** a client sends the `products` query and products exist
- **THEN** the system returns the full list of products

#### Scenario: No products
- **WHEN** a client sends the `products` query and no products exist
- **THEN** the system returns an empty list

### Requirement: Get product by id
The system SHALL return a single product by id through the `product` query, and MUST report a not-found condition as a typed GraphQL error.

#### Scenario: Product exists
- **WHEN** a client sends `product(id)` for an existing id
- **THEN** the system returns the matching product

#### Scenario: Product does not exist
- **WHEN** a client sends `product(id)` for an unknown id
- **THEN** the system returns a GraphQL error with extension code `PRODUCT_NOT_FOUND`

### Requirement: Create product
The system SHALL create a product through the `createProduct` mutation, delegating validation to the domain, and MUST NOT place business logic in the resolver.

#### Scenario: Valid input
- **WHEN** a client sends `createProduct` with a non-empty name and a price greater than zero
- **THEN** the system persists the product and returns it with a generated id and createdAt

#### Scenario: Invalid input
- **WHEN** a client sends `createProduct` with an empty name
- **THEN** the system returns a GraphQL error with extension code `INVALID_NAME` and persists nothing

### Requirement: Update product
The system SHALL update the name and/or price of a product through the `updateProduct` mutation.

#### Scenario: Existing product
- **WHEN** a client sends `updateProduct(id, input)` for an existing product with a new name
- **THEN** the system applies the change and returns the updated product

#### Scenario: Unknown product
- **WHEN** a client sends `updateProduct(id, input)` for an unknown id
- **THEN** the system returns a GraphQL error with extension code `PRODUCT_NOT_FOUND`

### Requirement: Delete product
The system SHALL delete a product through the `deleteProduct` mutation.

#### Scenario: Existing product
- **WHEN** a client sends `deleteProduct(id)` for an existing product
- **THEN** the system removes it and returns true

#### Scenario: Unknown product
- **WHEN** a client sends `deleteProduct(id)` for an unknown id
- **THEN** the system returns a GraphQL error with extension code `PRODUCT_NOT_FOUND`

### Requirement: Domain error mapping
The system SHALL map typed domain errors to GraphQL errors that carry an `extensions.code` and a `message`.

#### Scenario: Known domain error
- **WHEN** a resolver receives a typed domain error such as `ErrInvalidPrice`
- **THEN** the GraphQL error exposes the corresponding code, for example `INVALID_PRICE`

#### Scenario: Unknown error
- **WHEN** a resolver receives an error that is not a known domain error
- **THEN** the GraphQL error exposes the code `INTERNAL_ERROR`
