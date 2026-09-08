# ADR 006: Consistency Model

## Context

The simple implementation writes memberships directly to Redis through the REST API.

The production-ready architecture processes memberships asynchronously through Kafka, so there can be a short delay before a new event appears in the estimate result.

## Decision

Use read-after-write consistency for the simple implementation.

When the REST request succeeds, the membership has already been written to Redis, so the next estimate request can see the update.

Use eventual consistency for the production-ready architecture.

Kafka events are processed asynchronously, so the estimate result may briefly lag behind the latest events.

## Why This Approach?

The simple implementation does not need a distributed consistency model because it uses a single Redis instance and synchronous writes.

For the production-ready design, eventual consistency is a reasonable trade-off for better scalability, buffering, and reliability.

## Consequences

The simple version gives predictable results after a successful write.

The production-ready version may have a small delay between an event being accepted and the count being updated.

This delay should be monitored through Kafka consumer lag and service metrics.
