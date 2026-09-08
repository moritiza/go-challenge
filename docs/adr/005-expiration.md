# ADR 005: Membership Expiration

## Context

Each `(user_id, segment)` membership should stay active for 14 days.

We need to know if a membership is still active when calculating the segment count.

## Decision

Store the expiration timestamp as the score of a Redis Sorted Set.

For example:

`segment:sports`

- `user-101 -> expiration_timestamp`
- `user-202 -> expiration_timestamp`

When a membership is received, its expiration time is set to 14 days from the event processing time.

If the same `(user_id, segment)` is received again, the expiration time is updated and the 14-day period starts again.

When calculating the count, only memberships with an expiration time greater than the current time are counted.

Expired members can be removed later by a background cleanup process.

## Why This Approach?

- Redis Sorted Sets fit this use case well.
- Counting active memberships is simple and fast.
- We do not need to delete every membership exactly when it expires.
- The same data model works with both the simple Redis setup and Redis Cluster.

## Consequences

The count stays correct even if expired members have not been deleted yet.

Redis memory can grow because expired members may remain for some time, so a cleanup process should periodically remove them.

In the production-ready architecture, Kafka keeps the original events. Redis can therefore be rebuilt by replaying events and applying the same 14-day expiration rule.
