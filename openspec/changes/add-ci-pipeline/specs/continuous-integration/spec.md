## ADDED Requirements

### Requirement: Continuous integration pipeline
The system SHALL run an automated CI pipeline on every push and pull request to the `main` and `develop` branches, covering linting, vet, build, tests with the race detector and the Docker image build.

#### Scenario: Pull request checks
- **WHEN** a pull request targets `main` or `develop`
- **THEN** the pipeline runs lint, vet, build, race tests and the Docker build, and must pass for the change to be considered green

#### Scenario: Module tidiness
- **WHEN** the pipeline runs
- **THEN** it fails if `go.mod` or `go.sum` are not tidy
