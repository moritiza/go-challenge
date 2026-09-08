# ADR 007: Testing Strategy

## Context

Testing is important for this service because the main logic involves memberships, expiration, counting, and error handling.

The main focus is on unit tests that are fast, simple, and easy to maintain.

## Decision

Use unit tests as the main testing strategy.

Business logic should be tested without depending on Redis, Kafka, or HTTP.

We use interfaces and dependency injection so the service can be tested with a fake storage.

### Unit Tests

Important cases include:

- Store a new membership.
- Refresh an existing membership.
- Count active memberships.
- Ignore expired memberships.
- Handle invalid input.
- Handle storage errors.
- Handle different segments correctly.

### HTTP Handler Tests

HTTP handlers can be tested separately to verify:

- Request validation.
- Correct status codes.
- Response format.
- Error handling.

### Race Detection

Run the test suite with Go's race detector:

`go test -race ./...`

This helps catch data races in concurrent code.

## Testability

Business logic should not depend directly on Redis, Kafka, or HTTP.

Using interfaces and dependency injection keeps the code easy to test and keeps unit tests fast.

The current implementation focuses on unit tests. Integration tests can be added later if needed.

## Consequences

Most important business behavior can be tested without external services.

The test suite stays fast and simple while still covering the main success and failure cases.
