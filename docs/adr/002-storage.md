# ADR 002: Storage Choice

## Context

The Estimation Service needs to store user memberships and quickly count active users in each segment.

For the simple implementation, we want a storage solution that is fast, simple, and easy to test.

## Decision

Use Redis Sorted Sets for the simple implementation.

Each segment has its own Sorted Set:

- `member` = `user_id`
- `score` = expiration timestamp

This makes it easy to add or refresh a membership and count only active users.

The production-ready design uses Redis Cluster for higher traffic and availability. Kafka keeps the original events so Redis can be rebuilt if needed.

## Why Redis?

- Fast reads and writes.
- Sorted Sets fit the expiration and counting use case well.
- No need for complex queries.
- Easy to run for the simple implementation.
- Can scale to Redis Cluster in the production design.

## Alternatives

### PostgreSQL

PostgreSQL provides strong durability and rich queries, but it adds more storage and query complexity for this use case.

### MongoDB

MongoDB could store the memberships, but its document model does not give us a clear advantage for this access pattern.

## Consequences

The simple implementation stays small and fast, but Redis is not used as the durable source of events.

In the production-ready design, Kafka provides durable events and Redis acts as the fast serving store. This also allows Redis data to be rebuilt if needed.
