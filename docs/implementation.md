# Implementation Notes

> **Note:** The current implementation follows the [simple architecture](architectures/simple-architecture.md) and intentionally favors the Go standard library, readability, and simplicity. For a production deployment at significantly higher scale, the [production-ready architecture](architectures/production-ready-architecture.md) should be used.

## Expiration and Cleanup

Membership expiration is part of the data model itself.

Each segment is stored as a Redis Sorted Set where:

- `member` is the `user_id`
- `score` is the membership expiration timestamp

Active users are counted directly using the expiration score, so expired memberships are excluded from estimates even if they have not yet been physically removed from Redis.

This means cleanup is an operational concern rather than a correctness requirement. A delayed cleanup does not change the result returned by the estimation API.

Expired memberships are removed by a background worker in bounded batches instead of loading or deleting all expired members at once.

The batch size is configurable and keeps memory usage bounded regardless of the total number of users in a segment. It also prevents a large cleanup operation from doing an unbounded amount of work in a single iteration.

The cleanup process is:

1. Read a bounded number of expired members from Redis.
2. Remove those members.
3. Repeat while expired members remain.

A Lua script could combine the read and delete operations, but it was intentionally not used here.

Cleanup does not participate in the correctness of `Estimate`, so atomic cleanup is not required. Keeping the operations separate makes the implementation easier to read, test, and maintain.

For this challenge, the trade-off favors simplicity over server-side scripting. In production, cleanup can be further optimized with techniques such as Lua scripting, pipelining, or parallel cleanup workers if monitoring shows that cleanup has become a bottleneck.

## HTTP Layer

The service uses Go's standard `net/http` package and `http.ServeMux`.

The API surface is small, so adding an external framework would add dependencies and abstraction without significant benefit for this implementation.

For the production architecture, membership ingestion should preferably move to Kafka consumers as described in the [production-ready architecture](architectures/production-ready-architecture.md).

If synchronous HTTP ingestion is still required at very high request rates, the HTTP layer should be benchmarked and optimized based on actual workload. A framework such as Fiber can be considered if it provides a measurable benefit.

## Logging

The implementation uses Go's standard `log/slog` package because it provides structured logging and JSON output without adding another dependency.

For high-throughput production workloads, a logger such as Zap can be considered if profiling shows logging overhead to be significant.

Production deployments should also separate log generation from log collection. The service can write structured logs to stdout while an agent or sidecar forwards them to the centralized logging system.

## Observability and Tracing

Observability was intentionally kept lightweight in this implementation to focus on the core challenge and avoid adding unnecessary dependencies.

A production version should use OpenTelemetry for application metrics and distributed tracing across HTTP, Kafka, Redis, and important business operations.

A production observability stack should monitor request latency and errors, ingestion throughput, Redis performance, cleanup activity, Kafka consumer lag, and resource usage.

Metrics and traces can then be exported to the organization's monitoring and tracing infrastructure, such as Prometheus/Grafana and an OpenTelemetry-compatible tracing backend.
