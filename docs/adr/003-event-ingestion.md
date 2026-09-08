# ADR 003: Event Ingestion

## Context

The Estimation Service needs to receive a large number of membership events from USS.

For the simple implementation, we want to keep ingestion synchronous and avoid extra infrastructure.

For the production-ready design, we need better scaling, buffering, retries, and reliable event processing.

## Decision

Use REST for ingestion in the simple implementation.

USS sends membership data directly to the Estimation Service using:

`POST /v1/memberships`

For the production-ready architecture, use Kafka for asynchronous event ingestion.

USS publishes events to the `user-segment-events` topic, and Estimation Service workers consume them as a consumer group.

## Why Kafka for Production?

- It can handle high event volume.
- It provides buffering when consumers are busy.
- Multiple workers can process events in parallel.
- Failed events can be retried.
- Events can be kept and replayed when needed.
- It reduces the dependency between USS and ES.

## Alternatives

### Direct REST

Simple and easy to implement, but the producer depends directly on ES availability and there is no built-in event buffering.

### RabbitMQ

RabbitMQ is good for message-based workloads, but Kafka fits better because we want durable events, replay, and high-throughput event processing.

## Consequences

The simple implementation stays small and easy to run.

The production-ready design adds Kafka infrastructure and more operational complexity, but gives us better scaling, reliability, and recovery.
