# TUrnUp Scheduler

Scheduler is a Go service responsible for pulling and storing TU event data
in Redis KV store for fast, and accessible, data access without needing to
interact with Towson's events APIs.

## Why?

TUrnUp's old architecture for pulling event data was based on a single Next.JS
server powered with tRPC. This architecture is fine for small scale use, but it
raised various issues in terms of latency and potential rate-limit concerns.

### Reason 1. Latency

Each request will fetch a list of events from both data sources, and map those
events to a standardized format. This process introduced a lot of latency since
we'd have to process dozens of events at once and return them to the user.

Obviously latency won't matter as much because this is simply an events aggregations
site, and we're not Google. But being fast is the new and cool thing to do, so
why not worry about it?

### Reason 2. Ratelimiting

These events APIs are unauthenticated and anonymous requests, so it's important to assume
that the APIs will behave defensively. Scheduler's goal is to minimize the amount
of requests it takes to retrieve events, by storing these lists of events in our own
memory cache and serving this data from cache.
