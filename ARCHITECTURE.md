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
│           ├── service_name.proto         ← gRPC + HTTP service definition (google.api.http annotations)
│           ├── service_name.pb.go         ← generated: messages & enums
│           ├── service_name_grpc.pb.go    ← generated: gRPC client/server stubs
│           ├── service_name.pb.gw.go      ← generated: gRPC-Gateway reverse proxy (HTTP→gRPC)
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
        │   ├── repository/
        │   │   └── user.go               ← DB projection models (db tags)
        │   └── gateway/
        │       └── submit_order.go        ← outbound call DTOs
        │
        ├── handler/                       ← inbound: validate → call service → respond
        │   ├── http/
        │   │   ├── http.go               ← gRPC-Gateway setup + non-proto endpoints
        │   │   └── create_order.go        ← example: non-proto endpoint (file upload, SSE)
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
    http.go                  — gRPC-Gateway registration + non-proto endpoints
    create_order.go          — example: non-proto endpoint (file upload, SSE)
  grpc/
    grpc.go                  — GrpcHandler struct + constructor
    find_driver_by_status.go — FindDriverByStatus RPC handler
  kafka/
    kafka.go                 — KafkaHandler struct + constructor
    user_registered.go       — consumes "user.registered" topic
```

**Proto is the single source of truth for all transports:**

| Transport | DTO source | How it works |
|-----------|-----------|--------------|
| HTTP      | `pbentity` (protobuf-generated) | `google.api.http` annotations in proto → gRPC-Gateway reverse proxy auto-translates HTTP/JSON ↔ gRPC |
| gRPC      | `pbentity` (protobuf-generated) | Proto definition IS the contract |
| Kafka     | `pbentity` (protobuf-generated) | Proto definition IS the contract |

**HTTP via gRPC-Gateway:**

HTTP endpoints are defined as `google.api.http` annotations on gRPC RPCs in the
proto file. `protoc-gen-grpc-gateway` generates a reverse proxy (`*.pb.gw.go`)
that translates REST/JSON requests into gRPC calls. This means:

- No manual HTTP DTOs (`entity/http` is eliminated).
- No hand-written HTTP route registration for proto-defined endpoints.
- Swagger/OpenAPI spec is auto-generated by `protoc-gen-openapiv2`.
- The gRPC handler implements the business logic; HTTP is just a transport facade.

`handler/http/http.go` has two responsibilities:
1. **gRPC-Gateway registration** — calls `pbentity.RegisterServiceNameServiceHandlerFromEndpoint()` to wire up the reverse proxy.
2. **Non-proto endpoints** — any route that cannot be expressed in proto (file uploads, WebSocket, SSE, health probes).

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
  repository/    — DB projection models (db tags)
    user.go      — User struct with `db:` tags
  gateway/       — Outbound call DTOs
    submit_order.go
```

There is **no** `entity/http/` sub-package. HTTP transport types are defined in
proto files and auto-generated — protobuf messages serve as HTTP DTOs via
gRPC-Gateway.

---

### `shared/`

A **separate Go module** (`shared/go.mod`) imported by all services.

```
shared/
  consts/     — package-level constants shared across services
  errorsx/    — sentinel errors (ErrInternal, ErrNotFound, …)
  proto/
    service_name/
      service_name.proto        — gRPC + HTTP service definition (google.api.http annotations)
      service_name.pb.go        — generated: messages & enums
      service_name_grpc.pb.go   — generated: gRPC client/server stubs
      service_name.pb.gw.go     — generated: gRPC-Gateway reverse proxy (HTTP→gRPC)
      kafka.proto               — Kafka event schemas
      kafka.pb.go               — generated
  types/      — shared value objects: enums with Scan/Value/ToPB/FromPB
```

**Proto as the canonical contract:**

All service endpoints — HTTP, gRPC, and Kafka — are defined in proto files.
HTTP endpoints use `google.api.http` annotations on gRPC RPCs, which
`protoc-gen-grpc-gateway` compiles into a reverse-proxy handler.
Swagger/OpenAPI is auto-generated by `protoc-gen-openapiv2`.

The gRPC-Gateway registration (calling `RegisterXxxHandlerFromEndpoint`)
lives in `handler/http/http.go` inside each service — not in `shared/proto/`.

**Enum type safety:**

`types/` enums use the **generated proto enum type** in `ToPB()`/`FromPB()`,
not raw `int32`. This makes conversions self-documenting and compile-safe:

```go
func (f FinanceInvoiceStatus) ToPB() pbentity.Status { ... }
func FinanceInvoiceStatusFromPB(status pbentity.Status) FinanceInvoiceStatus { ... }
```

**What belongs here:** technical primitives, transport-neutral value types, proto definitions.  
**What does NOT belong here:** domain logic, service-specific error types, business rules.

---

## Data Flow

### Inbound HTTP Request (via gRPC-Gateway)

```
HTTP client
    │  POST /v1/invoices  (JSON body)
    ▼
gRPC-Gateway  (service_name.pb.gw.go — auto-generated)
    │  JSON → pbentity.CreateRequestedRequest (automatic)
    │  proxies to gRPC server on localhost
    ▼
handler/grpc/create_requested.go
    │  validate pbentity.CreateRequestedRequest
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
handler/grpc/create_requested.go
    │  map → pbentity.InvoiceResponse
    ▼
gRPC-Gateway
    │  pbentity.InvoiceResponse → JSON (automatic)
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
       ├── pbentity        (shared/proto/service_name — gRPC-Gateway HTTP setup)
       ├── repository/*    depends on → entity/repository
       ├── gateway/*       depends on → entity/gateway, pbentity
       ├── service/*       depends on → entity/service
       │                   injects    → repository interfaces
       │                   injects    → gateway interfaces
       └── handler/*       depends on → pbentity
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
| 4 | **Layer-scoped DTOs** | Each layer owns its types. DB models carry `db` tags; the domain model is tag-free. HTTP and gRPC share proto-generated types — no hand-maintained HTTP DTOs. |
| 5 | **Consistent transaction model** | `PgAtomicRepository` (Start/Finish/Abort) provides a single, reusable transaction abstraction. `transactionalOperation` in service eliminates boilerplate across use cases. |
| 6 | **Proto as the single source of truth** | For **all transports** (HTTP, gRPC, Kafka), protobuf definitions live in `shared/proto/` and are the canonical contracts. HTTP endpoints are declared via `google.api.http` annotations; gRPC-Gateway generates the REST layer. Swagger/OpenAPI is auto-generated — no hand-maintained API docs. |
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

---

## Guides

---

### Working with Enum Types

Enums live in `shared/types/` and bridge three worlds: **Go domain logic**, **database storage**, and **protobuf transport**. Every enum follows the same structure.

#### 1. Define the proto enum

Every enum starts in a `.proto` file inside `shared/proto/service_name/`.
`protoc` generates the Go type and named constants automatically.

```proto
// shared/proto/service_name/service_name.proto
enum Status {
  INVOICE_STATUS_UNSPECIFIED = 0;
  INVOICE_STATUS_REQUESTED = 2;
  INVOICE_STATUS_INVOICED = 3;
  INVOICE_STATUS_PAID = 4;
  INVOICE_STATUS_PARTIALLY_PAID = 5;
  INVOICE_STATUS_VOIDED = 6;
}
```

After generation, the Go package exposes:
- Type: `pbentity.Status`
- Constants: `pbentity.Status_INVOICE_STATUS_REQUESTED`, etc.

#### 2. Create the Go enum in `shared/types/`

Each enum file has a fixed set of components:

```go
package types

import (
    "database/sql/driver"
    pbentity "awesomeProject/shared/proto/service_name"
)

// 1) Go type — used in domain/service layer (no proto dependency there).
type FinanceInvoiceStatus int

// 2) iota constants — internal representation.
const (
    FinanceInvoiceStatusUnspecified FinanceInvoiceStatus = iota
    FinanceInvoiceStatusDraft
    FinanceInvoiceStatusRequested
    // ...
)

// 3) String constants — match the database representation.
const (
    FinanceInvoiceStatusUnspecifiedString = "unspecified"
    FinanceInvoiceStatusDraftString       = "draft"
    FinanceInvoiceStatusRequestedString   = "requested"
    // ...
)

// 4) Mapping tables — four maps for bidirectional conversion.
var (
    // Go → Proto: uses typed proto enum constants, NOT raw int32.
    financeInvoiceStatusToPBMap = map[FinanceInvoiceStatus]pbentity.Status{
        FinanceInvoiceStatusRequested: pbentity.Status_INVOICE_STATUS_REQUESTED,
        FinanceInvoiceStatusInvoiced:  pbentity.Status_INVOICE_STATUS_INVOICED,
        // ...
    }

    // Proto → Go
    pbToFinanceInvoiceStatusMap = map[pbentity.Status]FinanceInvoiceStatus{
        pbentity.Status_INVOICE_STATUS_REQUESTED: FinanceInvoiceStatusRequested,
        pbentity.Status_INVOICE_STATUS_INVOICED:  FinanceInvoiceStatusInvoiced,
        // ...
    }

    // Go → String (for DB and JSON)
    financeInvoiceStatusToStringMap = map[FinanceInvoiceStatus]string{ /* ... */ }

    // String → Go (from DB and JSON)
    stringToFinanceInvoiceStatusMap = map[string]FinanceInvoiceStatus{ /* ... */ }
)
```

#### 3. Required methods

Every enum must implement these methods:

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Scan` | `(f *T) Scan(value interface{}) error` | `database/sql.Scanner` — reads from DB |
| `Value` | `(f T) Value() (driver.Value, error)` | `database/sql/driver.Valuer` — writes to DB |
| `String` | `(f T) String() string` | Human-readable / DB string form |
| `ToPB` | `(f T) ToPB() pbentity.ProtoEnum` | Converts to **typed** proto enum |
| `FromPB` | `TFromPB(status pbentity.ProtoEnum) T` | Converts from **typed** proto enum |
| `FromString` | `TFromString(s string) T` | Parses from DB / JSON string |

**Critical rule:** `ToPB()` returns the **generated proto enum type** (e.g. `pbentity.Status`),
not `int32`. `FromPB()` accepts the proto enum type. This gives compile-time safety —
if a proto value is renamed or removed, the code stops compiling instead of silently
mapping to a wrong number.

```go
// CORRECT — type-safe, self-documenting:
func (f FinanceInvoiceStatus) ToPB() pbentity.Status {
    if pb, exists := financeInvoiceStatusToPBMap[f]; exists {
        return pb
    }
    return pbentity.Status_INVOICE_STATUS_UNSPECIFIED
}

func FinanceInvoiceStatusFromPB(status pbentity.Status) FinanceInvoiceStatus {
    if s, exists := pbToFinanceInvoiceStatusMap[status]; exists {
        return s
    }
    return FinanceInvoiceStatusUnspecified
}
```

```go
// WRONG — magic numbers, no compile-time safety:
func (f FinanceInvoiceStatus) ToPB() int32 {
    // ...
    return 2  // what is 2? breaks silently if proto renumbers
}
```

#### 4. Unspecified / zero value

Both the Go enum (`iota` starts at 0) and the proto enum (first value `= 0`) use
zero as the "unspecified" sentinel. All `FromX` functions return the `Unspecified`
variant when the input doesn't match any known value — never panic.

#### 5. Adding a new enum — checklist

1. Add `enum NewStatus { ... }` to the relevant `.proto` file.
2. Run `protoc` to regenerate `.pb.go`.
3. Create `shared/types/new_status.go` following the template above.
4. Implement all 6 methods + 4 maps.
5. Use `NewStatus` in domain models; convert at layer boundaries with `ToPB()`/`FromPB()`.

---

### Generating HTTP from gRPC (gRPC-Gateway)

HTTP endpoints are not hand-written. They are derived from proto definitions
via [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway).

#### How it works

```
service_name.proto          (you write: RPC + google.api.http annotation)
        │
        ▼  protoc
service_name.pb.go          (generated: messages, enums)
service_name_grpc.pb.go     (generated: gRPC stubs)
service_name.pb.gw.go       (generated: HTTP reverse proxy)
service_name.swagger.json   (generated: OpenAPI spec)
```

The generated `*.pb.gw.go` is a Go HTTP handler that:
1. Accepts an HTTP/JSON request.
2. Deserializes JSON into the protobuf request message.
3. Forwards it to the gRPC server on localhost.
4. Serializes the protobuf response back to JSON.

#### Step 1 — Annotate RPCs in proto

Add `google.api.http` options to each RPC that should be accessible over HTTP.
RPCs without the annotation remain gRPC-only.

```proto
import "google/api/annotations.proto";

service ServiceNameService {
  // HTTP + gRPC: has annotation → exposed as POST /v1/invoices.
  rpc CreateRequested(CreateRequestedRequest) returns (InvoiceResponse) {
    option (google.api.http) = {
      post: "/v1/invoices"
      body: "*"
    };
  }

  // gRPC only: no annotation → not exposed over HTTP.
  rpc FindDriverByStatus(FindDriverByStatusRequest) returns (FindDriverByStatusResponse);
}
```

Common HTTP method patterns:

```proto
// POST with JSON body
option (google.api.http) = { post: "/v1/resources" body: "*" };

// GET with path parameter
option (google.api.http) = { get: "/v1/resources/{id}" };

// PUT with JSON body
option (google.api.http) = { put: "/v1/resources/{id}" body: "*" };

// DELETE
option (google.api.http) = { delete: "/v1/resources/{id}" };

// GET with query parameters (all message fields become ?key=value)
option (google.api.http) = { get: "/v1/resources" };
```

#### Step 2 — Generate code

```bash
protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
  --openapiv2_out=. \
  -I . \
  -I $(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@latest/third_party/googleapis \
  shared/proto/service_name/service_name.proto
```

This produces:
- `service_name.pb.go` — messages & enums
- `service_name_grpc.pb.go` — gRPC client/server interfaces
- `service_name.pb.gw.go` — HTTP reverse proxy handler
- `service_name.swagger.json` — OpenAPI/Swagger spec

#### Step 3 — Register the gateway in `handler/http/http.go`

```go
package http

import (
    "context"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pbentity "awesomeProject/shared/proto/service_name"
)

type HttpHandler struct {
    gatewayMux *runtime.ServeMux
}

func New(ctx context.Context, grpcAddr string) (*HttpHandler, error) {
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    if err := pbentity.RegisterServiceNameServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
        return nil, err
    }

    return &HttpHandler{gatewayMux: mux}, nil
}

func (h *HttpHandler) Handler() http.Handler {
    return h.gatewayMux
}
```

#### Step 4 — Wire in `app/run.go`

```go
httpHandler, err := handlerhttp.New(ctx, cfg.GRPC.Addr)
if err != nil {
    return fmt.Errorf("http handler init: %w", err)
}

httpServer := &http.Server{
    Addr:    ":8080",
    Handler: httpHandler.Handler(),
}
```

#### Request lifecycle

```
HTTP client  →  POST /v1/invoices {"organization_id":1, "note":"test"}
                        │
                        ▼
              gRPC-Gateway (*.pb.gw.go)
                JSON → pbentity.CreateRequestedRequest
                        │
                        ▼  (localhost gRPC call)
              handler/grpc/create_requested.go
                validates, maps to domain model, calls service
                        │
                        ▼
              service layer → repository → ...
                        │
                        ▼
              handler/grpc/create_requested.go
                maps result → pbentity.InvoiceResponse
                        │
                        ▼
              gRPC-Gateway
                pbentity.InvoiceResponse → JSON
                        │
                        ▼
HTTP client  ←  200 OK {"invoice": {...}}
```

#### Required protoc plugins

| Plugin | Installed via | Generates |
|--------|---------------|-----------|
| `protoc-gen-go` | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` | `*.pb.go` |
| `protoc-gen-go-grpc` | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` | `*_grpc.pb.go` |
| `protoc-gen-grpc-gateway` | `go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest` | `*.pb.gw.go` |
| `protoc-gen-openapiv2` | `go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest` | `*.swagger.json` |

---

### Working with `transactionalOperation`

`transactionalOperation` is a reusable helper in the service layer that wraps
a function in a Postgres transaction with automatic commit/rollback.

#### How it works

```
service.transactionalOperation(ctx, func(tx) error {
    ┌──────────────────────────────────────────┐
    │  repo.Start(ctx)         → begin TX      │
    │                                          │
    │  txFunc(tx)              → your logic    │
    │    ├── repo.Create(ctx, tx, ...)         │
    │    ├── repo.Update(ctx, tx, ...)         │
    │    └── return err                        │
    │                                          │
    │  if err != nil:                          │
    │    repo.Abort(tx)        → ROLLBACK      │
    │  else:                                   │
    │    repo.Finish(tx)       → COMMIT        │
    └──────────────────────────────────────────┘
})
```

The three primitives come from `PgAtomicRepository`:

| Method | SQL | When |
|--------|-----|------|
| `Start(ctx)` | `BEGIN` | Opens a new transaction, returns `*sqlx.Tx` |
| `Finish(tx)` | `COMMIT` | Called when `txFunc` returns `nil` |
| `Abort(tx)` | `ROLLBACK` | Called in `defer` when `txFunc` returns an error; ignores `sql.ErrTxDone` |

#### Usage pattern

```go
func (s *service) UpdateUser(ctx context.Context, id string, user *serviceEntity.User) (*serviceEntity.User, error) {
    var result *serviceEntity.User

    err := s.transactionalOperation(ctx, func(tx *sqlx.Tx) error {
        // All repository calls inside this closure share the same transaction.
        // Pass `tx` to every repository method that needs to participate.

        existing, err := s.repo.Pg.FetchById(ctx, tx, id)
        if err != nil {
            return err  // → triggers Abort (ROLLBACK)
        }

        existing.Username = user.Username
        existing.Email = user.Email

        updated, err := s.repo.Pg.Update(ctx, tx, existing)
        if err != nil {
            return err  // → triggers Abort (ROLLBACK)
        }

        result = updated
        return nil  // → triggers Finish (COMMIT)
    })
    if err != nil {
        return nil, err
    }

    return result, nil
}
```

#### Rules

1. **Every repository method inside the closure must accept `*sqlx.Tx`** —
   this ensures all operations run in the same transaction.

2. **Return early on error** — any non-nil error returned from `txFunc`
   triggers `Abort` (ROLLBACK) via the deferred function.

3. **Do not call `Finish` or `Abort` manually** — `transactionalOperation`
   handles both automatically. Calling them yourself will cause double-commit
   or double-rollback.

4. **Keep the closure short** — long-running logic inside a transaction holds
   row locks and increases the risk of deadlocks. Fetch/validate outside the
   transaction when possible; only put writes inside.

5. **Capture results via closure variables** — the closure returns only `error`.
   Use a variable declared before the closure (like `result` above) to pass
   data back to the calling function.

6. **One transaction = one storage engine** — `transactionalOperation` wraps
   Postgres only. If you need cross-storage atomicity (e.g. Postgres + Kafka),
   use the outbox pattern or saga instead.

#### When to use vs. when to skip

| Scenario | Use `transactionalOperation`? |
|----------|-------------------------------|
| Multiple writes that must be atomic | Yes |
| Single INSERT/UPDATE | No — direct repo call is simpler |
| Read-only queries | No — no transaction needed |
| Writes to Postgres + publish to Kafka | Yes for PG writes; Kafka publish goes in `gateway/kafka/` after commit |
| Multi-step saga across services | No — use an explicit saga/outbox pattern |

