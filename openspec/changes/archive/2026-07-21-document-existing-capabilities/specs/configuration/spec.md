## ADDED Requirements

### Requirement: Environment configuration with defaults
The system SHALL load its configuration from the environment, applying sensible defaults for every value when the variable is absent.

#### Scenario: Defaults
- **WHEN** no configuration variables are set
- **THEN** each field takes its default value

#### Scenario: Overrides
- **WHEN** configuration variables are set
- **THEN** each field takes the value from the environment

#### Scenario: Invalid boolean
- **WHEN** a boolean variable holds a non-boolean value
- **THEN** the field falls back to its default instead of failing
