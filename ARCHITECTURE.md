# Microservice Architecture

## Overview

This project follows a **Clean Architecture** (Ports & Adapters) pattern for each microservice.
The core principle: **business logic has zero knowledge of transport, storage, or external services**.
Dependencies always point inward — toward the domain, never outward.

```
 External World
       │
       ▼
  ┌─────────┐      ┌─────────┐
  │ handler │      │ gateway │
  └────┬────┘      └────┬────┘
       │                │
       ▼                ▼
  ┌──────────────────────────┐
  │         service          │  ← business logic lives here only
  └──────────┬───────────────┘
             │
             ▼
       ┌──────────┐
       │repository│
       └──────────┘
```

---

## Full Project Tree

```
awesomeProject/
│
├── shared/                                ← separate Go module (awesomeProject/shared)
│   ├── go.mod                             ← module root for shared; imported by every service
│   ├── consts/
│   │   └── consts.go
│   ├── errorsx/
│   │   └── errors.go
│   ├── types/
│   │   ├── example.go                     ← shared value types (enums, custom scalars)
│   │   └── example2.go
│   └── proto/
│       └── service_name/
│           ├── service_name.proto         ← gRPC service definition
│           ├── service_name.pb.go         ← generated
│           ├── kafka.proto                ← Kafka event schemas
│           └── kafka.pb.go               ← generated
│
└── services/
    └── service_name/                      ← one directory per microservice
        ├── go.mod                         ← module root for this service (awesomeProject/services/service_name)
        │
        ├── cmd/
        │   └── main.go                    ← binary entry point; calls app.Run()
        │
        ├── config/
        │   └── config.go                  ← Config struct + New()
        │
        ├── app/                           ← composition root; infrastructure only
        │   ├── run.go                     ← wires deps, starts servers, graceful shutdown
        │   ├── db/
        │   │   ├── pg/
        │   │   │   ├── pg.go              ← Postgres client; migrations run inside New()
        │   │   │   └── migrations/
        │   │   ├── mongo/
        │   │   │   ├── mongo.go           ← Mongo client; migrations run inside New()
        │   │   │   └── migrations/
        │   │   └── redis/
        │   │       └── redis.go           ← Redis client
        │   ├── grpc/
        │   │   ├── server.go              ← gRPC server setup & handler registration
        │   │   └── client.go             ← gRPC client connection
        │   └── kafka/
        │       ├── publisher.go           ← Kafka producer setup
        │       └── consumer.go           ← Kafka consumer setup
        │
        ├── entity/                        ← data types scoped to each layer
        │   ├── service/
        │   │   └── user.go               ← domain models (pure Go structs, no tags)
        │   ├── http/                      ← [package: httpentity]
        │   │   └── user.go               ← HTTP DTOs (json tags); used only by handler/http/
        │   ├── repository/
        │   │   └── user.go               ← DB projection models (db tags)
        │   └── gateway/
        │       └── submit_order.go        ← outbound call DTOs
        │
        ├── handler/                       ← inbound: validate → call service → respond
        │   ├── http/
        │   │   ├── http.go               ← HttpHandler struct + New()
        │   │   └── create_order.go        ← POST /orders  (uses httpentity)
        │   ├── grpc/
        │   │   ├── grpc.go               ← GrpcHandler struct + New()
        │   │   └── find_driver_by_status.go  ← RPC handler (uses pbentity)
        │   └── kafka/
        │       ├── kafka.go              ← KafkaHandler struct + New()
        │       └── user_registered.go    ← consumes topic (uses pbentity)
        │
        ├── gateway/                       ← outbound: wraps all external calls
        │   ├── gateway.go                ← Gateway aggregate struct
        │   ├── grpc/
        │   │   ├── grpc.go              ← GrpcGateway interface + struct
        │   │   └── get_user.go           ← gRPC call to another service
        │   ├── http/
        │   │   ├── http.go              ← HttpGateway interface + struct
        │   │   └── collect_currency_rates.go  ← 3rd-party REST call
        │   └── kafka/
        │       ├── kafka.go             ← KafkaGateway interface + struct
        │       └── user_added.go         ← publish event to Kafka
        │
        ├── service/                       ← business logic only; one file per use case
        │   ├── service.go                ← Service interface + private struct + New()
        │   ├── get_user.go
        │   ├── update_user.go
        │   └── transactional_operation.go ← reusable tx helper (Start/txFunc/Finish/Abort)
        │
        └── repository/                    ← data access; one sub-package per storage engine
            ├── repository.go             ← Repository aggregate struct (Pg + Mongo + Redis)
            ├── pg/
            │   ├── pg.go                 ← PgRepository interface
            │   ├── pg_atomic.go          ← PgAtomicRepository: Start / Finish / Abort
            │   ├── create.go
            │   └── fetch_by_id.go
            ├── mongo/
            │   ├── mongo.go             ← MongoRepository interface
            │   └── create.go
            └── redis/
                ├── redis.go             ← RedisRepository interface
                ├── get.go
                └── set.go
```

---

## Module Structure

The project uses a **multi-module monorepo**: each service and the shared library are
independent Go modules. This means they version independently, have their own
`go.sum`, and a change in one does not force a rebuild of the others unless the
dependency is explicitly updated.

```
awesomeProject/
├── shared/
│   └── go.mod    ← module "awesomeProject/shared"
│
└── services/
    ├── service_name/
    │   └── go.mod    ← module "awesomeProject/services/service_name"
    │                    requires: awesomeProject/shared vX.Y.Z
    ├── service_two/
    │   └── go.mod    ← module "awesomeProject/services/service_two"
    └── ...
```

**Rules:**
- `shared/go.mod` has **no** dependencies on any `services/*` module — ever.
- Each service's `go.mod` declares `awesomeProject/shared` as a `require`.
- During local development use a `replace` directive to point at the local `shared/` path
  instead of a registry version:
  ```
  replace awesomeProject/shared => ../../shared
  ```
- Services **never** import each other directly. Cross-service communication is
  always over the network (gRPC / Kafka / HTTP).

---

## Layer-by-Layer Reference

---

### `cmd/`

**Responsibility:** Binary entry point. Nothing more.

Reads the config path from CLI arguments and calls `app.Run()`.
Must contain zero business logic.

```
cmd/
  main.go
```

---

### `config/`

**Responsibility:** Defines the `Config` struct and its constructor.

Loaded once at startup by `app/run.go`. Passed down as explicit dependencies —
never accessed as a global singleton.

```
config/
  config.go     — Config struct + New()
```

---

### `app/`

**Responsibility:** Infrastructure initialization, dependency wiring, and process lifecycle.

This is the **composition root** of the service — the only place where all concrete
implementations are instantiated and wired together. Nothing outside `app/` knows
about concrete types; everything else depends on interfaces.

```
app/
  run.go              — wires all dependencies, starts servers, handles shutdown
  db/
    pg/
      pg.go           — creates Postgres client; runs migrations inside New()
      migrations/     — SQL migration files
    mongo/
      mongo.go        — creates Mongo client; runs migrations inside New()
      migrations/
    redis/
      redis.go        — creates Redis client
  grpc/
    server.go         — gRPC server setup & registration
    client.go         — gRPC client connection setup
  kafka/
    publisher.go      — Kafka producer setup
    consumer.go       — Kafka consumer setup
```

**Key rules:**
- `New()` for every DB client runs migrations before returning. If migrations fail, the service does not start.
- `run.go` owns graceful shutdown: catches `SIGTERM`/`SIGINT`, drains HTTP and gRPC servers, flushes the Kafka producer.
- Structured logging (`log/slog`) is initialized here and injected as a dependency — never used as a global.

---

### `handler/`

**Responsibility:** Receives inbound requests, validates input, delegates to `service`, returns a response.

Handlers know nothing about storage or external services. Their only dependency is
the `service.Service` interface.

```
handler/
  http/
    http.go                  — HttpHandler struct + constructor
    create_order.go          — POST /orders handler
  grpc/
    grpc.go                  — GrpcHandler struct + constructor
    find_driver_by_status.go — FindDriverByStatus RPC handler
  kafka/
    kafka.go                 — KafkaHandler struct + constructor
    user_registered.go       — consumes "user.registered" topic
```

**Type sources by transport:**

| Transport | DTO source | Reason |
|-----------|-----------|--------|
| HTTP      | `entity/http` (`httpentity`) | Custom JSON contract, owned by this service |
| gRPC      | `pbentity` (protobuf-generated) | Proto definition IS the contract |
| Kafka     | `pbentity` (protobuf-generated) | Proto definition IS the contract |

gRPC and Kafka handlers do **not** use `entity/http` — protobuf types are the
wire format and the type at the same time, no intermediate DTO needed.

---

### `gateway/`

**Responsibility:** All outbound calls — to other internal services and to third-party APIs.

```
gateway/
  gateway.go           — Gateway aggregate struct (HttpGateway + GrpcGateway + KafkaGateway)
  grpc/
    grpc.go            — GrpcGateway interface + grpcGateway struct
    get_user.go        — calls another service via gRPC
  http/
    http.go            — HttpGateway interface + httpGateway struct
    collect_currency_rates.go — calls a third-party REST API
  kafka/
    kafka.go           — KafkaGateway interface + kafkaGateway struct
    user_added.go      — publishes "user.added" event
```

All gateway methods accept `context.Context` as the first argument to propagate
deadlines and cancellation across service boundaries.

DTOs for gateway calls live in `entity/gateway/`.

---

### `service/`

**Responsibility:** Business logic. The core of the service.

`service.go` defines the `Service` interface and the private `service` struct.
Each use case gets its own file.

```
service/
  service.go                  — Service interface + struct + New()
  get_user.go                 — GetUser use case
  update_user.go              — UpdateUser use case
  transactional_operation.go  — reusable transaction helper (Start / txFunc / Finish / Abort)
```

**Rules:**
- Depends only on `repository` interfaces and `gateway` interfaces — never on concrete types.
- Domain models come from `entity/service/`.
- No JSON tags, no `db` tags, no proto types — ever.
- `transactionalOperation` wraps repository's `Start`/`Finish`/`Abort` so individual
  use cases don't repeat the boilerplate.

---

### `repository/`

**Responsibility:** All data access. Abstracts storage technology behind interfaces.

```
repository/
  repository.go       — Repository aggregate struct (Pg + Mongo + Redis)
  pg/
    pg.go             — PgRepository interface
    pg_atomic.go      — PgAtomicRepository interface: Start / Finish / Abort
    create.go
    fetch_by_id.go
  mongo/
    mongo.go          — MongoRepository interface
    create.go
  redis/
    redis.go          — RedisRepository interface
    get.go
    set.go
```

**Transaction model (Postgres):**

```
PgRepository
  └── embeds PgAtomicRepository
        Start(ctx) → *sqlx.Tx
        Finish(tx) → commits
        Abort(tx)  → rolls back (ignores ErrTxDone)
```

`NewAtomic` returns `PgAtomicRepository` (interface) — not a pointer to an interface.

DB models live in `entity/repository/`.

---

### `entity/`

**Responsibility:** Layer-scoped data types. Each sub-package owns the types
for exactly one layer. If needed may include mappings between layers (try to avoid circular dependencies).

```
entity/
  service/       — domain models (no tags, pure Go structs)
    user.go      — User struct
  http/          — HTTP transport DTOs (JSON tags)   [package: httpentity]
    user.go      — CreateUserRequest, CreateUserResponse, ...
  repository/    — DB projection models (db tags)
    user.go      — User struct with `db:` tags
  gateway/       — Outbound call DTOs
    submit_order.go
```

---

### `shared/`

A **separate Go module** (`shared/go.mod`) imported by all services.

```
shared/
  consts/     — package-level constants shared across services
  errorsx/    — sentinel errors (ErrInternal, ErrNotFound, …)
  proto/
    service_name/
      service_name.proto   — gRPC service definition
      service_name.pb.go   — generated
      kafka.proto          — Kafka event schemas
      kafka.pb.go          — generated
  types/      — shared value objects: enums with Scan/Value/ToPB/FromPB
```

**What belongs here:** technical primitives, transport-neutral value types, proto definitions.  
**What does NOT belong here:** domain logic, service-specific error types, business rules.

---

## Data Flow

### Inbound HTTP Request

```
HTTP client
    │  POST /orders  (JSON body)
    ▼
handler/http/create_order.go
    │  decode → httpentity.CreateOrderRequest
    │  validate
    │  map → entity/service model
    ▼
service/create_order.go
    │  business rules
    │  calls repository and/or gateway
    ▼
repository/pg/create.go          gateway/grpc/get_user.go
    │  maps domain → entity/repository     │  maps domain → proto
    │  SQL INSERT with *sqlx.Tx            │  gRPC call to UserService
    ◄──────────────────────────────────────┘
    │  return entity/service model
    ▼
handler/http/create_order.go
    │  map → httpentity.CreateOrderResponse
    ▼
HTTP client  (JSON response)
```

### Inbound Kafka Message

```
Kafka broker  ("user.registered" topic)
    │  protobuf bytes
    ▼
handler/kafka/user_registered.go
    │  unmarshal → pbentity.UserRegisteredRequest
    │  validate / ensure idempotency
    │  map → entity/service model
    ▼
service/  (same path as above)
```

### Outbound Event

```
service/update_user.go
    │
    ▼
gateway/kafka/user_added.go
    │  maps → pbentity.UserAddedEvent
    │  producer.Publish("user.added", bytes)
    ▼
Kafka broker
```

---

## Dependency Graph

```
cmd
 └── app/run.go  (composition root)
       ├── config
       ├── app/db/*        (infrastructure)
       ├── app/grpc/*      (infrastructure)
       ├── app/kafka/*     (infrastructure)
       ├── repository/*    depends on → entity/repository
       ├── gateway/*       depends on → entity/gateway, pbentity
       ├── service/*       depends on → entity/service
       │                   injects    → repository interfaces
       │                   injects    → gateway interfaces
       └── handler/*       depends on → entity/http, pbentity
                           injects    → service.Service interface
```

**One-way rule:** every arrow points inward. `service` never imports `handler` or `gateway` packages directly.

---

## Pros & Cons

### ✅ Pros

| # | Advantage | Details |
|---|-----------|---------|
| 1 | **Transport independence** | Business logic is identical regardless of whether the request arrives over HTTP, gRPC, or Kafka. Adding a new transport means adding a new `handler/` sub-package — nothing else changes. |
| 2 | **Storage independence** | Swapping Postgres for another DB requires rewriting only `repository/pg/`. The service layer is untouched. |
| 3 | **Explicit composition root** | All wiring happens in one place (`app/run.go`). Tracing the dependency graph is trivial. No hidden globals, no `init()` surprises. |
| 4 | **Layer-scoped DTOs** | Each layer owns its types. HTTP handlers carry JSON tags; DB models carry `db` tags; the domain model is tag-free. Changes to one layer's schema do not cascade into others. |
| 5 | **Consistent transaction model** | `PgAtomicRepository` (Start/Finish/Abort) provides a single, reusable transaction abstraction. `transactionalOperation` in service eliminates boilerplate across use cases. |
| 6 | **Proto as the single source of truth** | For gRPC and Kafka, protobuf definitions live in `shared/proto/` and are the canonical contracts. No duplication between "our struct" and the wire format. |
| 7 | **Graceful shutdown built in** | `app/run.go` handles `SIGTERM`/`SIGINT`, drains in-flight requests, and flushes the Kafka producer before exit. Safe for containerized environments. |
| 8 | **Migrations are atomic with startup** | DB clients run migrations inside `New()`. If migrations fail, the process exits before accepting any traffic — no half-migrated state. |

---

### ❌ Cons & Known Trade-offs

| # | Problem | Impact |
|---|---------|--------|
| 1 | **Repository & Gateway God Objects** | `repository.Repository` and `gateway.Gateway` bundle all storage backends / transports into one struct. The service receives the entire toolbox instead of only what it needs. Violates ISP; makes mocking expensive. | 
| 2 | **`entity/service` must not import `entity/repository`** | The current template still shows `ToRepository()` on the domain model. This inverts the dependency. Mappers belong in `repository/`, not in the domain entity. |
| 3 | **`shared/types` couples services** | `FinanceInvoiceStatus`, `FeeFreeStatus` with `ToPB()`/`Scan()` in the shared module create compile-time coupling between otherwise independent services. A change to one type forces a rebuild of every service. |
| 4 | **Package name collisions** | `entity/service` (package `service`) and `service/` (package `service`) share the same package name. Any file importing both requires aliases, increasing cognitive load. |
| 5 | **No observability beyond logging** | The template bootstraps `slog` but has no OpenTelemetry tracing or Prometheus metrics hooks. In a distributed system, traces and metrics are as mandatory as logs. |
| 6 | **Kafka handlers lack idempotency contract** | Kafka guarantees at-least-once delivery. The template has no pattern for idempotency keys, deduplication, or Dead Letter Queue routing — this must be added before production use. |
| 7 | **`shared/errorsx` mixes domain boundaries** | `ErrUserNotFound` and `ErrPaymentNotFound` living in the same package implies those two domains are aware of each other. Each service should own its own domain errors. |
| 8 | **No context timeout enforcement in gateways** | Gateway methods accept `context.Context` but the template does not show wrapping calls with `context.WithTimeout`. Without it, a slow downstream blocks indefinitely. |

