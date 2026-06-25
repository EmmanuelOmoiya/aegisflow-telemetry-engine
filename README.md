# AegisFlow Telemetry Engine        

A high-throughput, concurrent event routing and policy enforcement engine built natively in Go for real-time endpoint compliance and sensor telemetry.

AegisFlow is engineered with zero external web framework dependencies, utilizing Go's native primitives to achieve low-overhead memory consumption and fast data distribution under heavy concurrency.

---

## Core Architecture & Data Flow

AegisFlow completely decouples network ingestion from data analysis loops, allowing the system to handle sudden traffic spikes without blocking application threads or droppping connections.

[Incoming Telemtry HTTP Request] -----------> internal/ingestion/handler.go - (Parses JSON & Increments Atomic Metrics) -----------> internal/queue/channel.go - (Drops Event into Bounded Go Channel) -----------> internal/worker/pool.go - (Static Goroutine Pool Pulls Concurrently) -----------> internal/policy/matcher.go - (Executes O(1) Policy In-Memory Lookup) ------> Matches? ----> [internal/alert/ ] [Log Execution Success] Doesn't Match ---> (Streams to Redis Pub/Sub)  

## Key Architectural Trade-offs

- **Standard Library Ingestion:** Built purely on `net/http` to bypass heavy routermiddleware allocations, maximizing raw request parsing speed.
- **Bounded In-Memory Backpressure:** The core pipelin cheannel is strictly bounded. If the ingestion queue saturates, the server safely triggers an active backpressure rejection (`HTTP 429 Too Many Requests`), protecting the host from Out-of-Memory (OOM) fatal drops.
- **Lock-Free Operational Metrics:** Live performance counters utilize low-level `sync/atomic` hardware primitives, bypassing traditional mutex contention over shared state variables.

---

## Folder Structure

```text
aegisflow/                  
|-- cmd/server/             # main.go - wires everything, handles SIGTERM
|-- internal/               
|    |-- ingestion/         # Payload types + HTTP handlers
|    |-- queue/             # Bounded channel queue wrapper
|    |-- worker/            # Concurrent goroutine worker pool
|    |-- policy/            # Rule loader + O(1) memory matcher
|    |-- alert/             # Redis Pub/Sub asynchronous publisher
|    |-- metrics/           # Atomic lock-free performance counters
|-- config/
|    |-- config.go          # Standard environment variable loader
|    |-- rules.yaml         # Default policy enforcement JSON/YAML rules
|-- scripts/
|    |-- seed.sh            # Benchmark load test invocation scripts
|-- Dockerfile              # Multi-stage container compilation -> scratch image
|-- docker-compose.yml      # Local validation composition (AegisFlow + Redis)
```