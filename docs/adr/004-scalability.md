# ADR 004: Scalability

## Context

The service needs to handle millions of users and a large number of membership events.

The simple implementation is intentionally small and runs with one service instance and one Redis instance.

The production-ready design needs to handle higher traffic without changing the main business logic.

## Decision

Use horizontal scaling for both the Estimation Service and the storage layer.

For the production-ready design:

- Run multiple Estimation Service workers in one Kafka consumer group.
- Use Kafka partitions to process events in parallel.
- Use Redis Cluster to spread data across multiple nodes.
- Shard Redis data by `segment`.

## Why This Approach?

- We can add more workers when event traffic grows.
- Kafka distributes the work between workers.
- Redis Cluster allows the storage layer to grow with the data.
- There is no need for a single large service instance.

### Alternatives

### Vertical Scaling

Adding more CPU and memory to a single Redis instance is simple, but it has limits and does not provide the same level of fault tolerance.

### Single Redis Instance

A single Redis instance is enough for the simple implementation, but it becomes a bottleneck as traffic and data grow.

## Consequences

The simple implementation stays easy to run and understand.

The production-ready design can scale horizontally with Kafka consumers and Redis Cluster, but it requires more infrastructure and operational work.

A very large or popular segment can still become a hot shard. If this becomes a problem, that segment can be split into smaller buckets.
