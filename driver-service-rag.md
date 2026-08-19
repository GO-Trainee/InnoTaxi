# Driver Service — RAG + LLM as the implementation foundation

**This track is optional.** Driver Service itself is required. RAG + LLM is not.

You can implement Driver Service by hand, the same way as User + Auth. That is a complete pass. Use this brief only if you *want* to try retrieval-augmented generation as the way you write the service.

**When:** after **User Service** and **Auth Service** are working, reviewed, and merged.  
**What:** implement Driver Service.  
**How (if you opt in):** do **not** write it from a blank chat with ChatGPT. Index *this* spec plus *your* User/Auth code, retrieve the relevant pieces, then generate and **review** the rest.

User + Auth are the house style. Driver Service must look like it belongs in the same monorepo, whether you used RAG or not.

## Why this optional track exists

A raw LLM does not know your project. It will happily invent Fiber + Postgres + GORM, a different error JSON, and a different folder layout — all “correct” in the abstract, all wrong *here*.

**RAG (Retrieval-Augmented Generation)** is the fix for that:

1. You turn the spec and the existing services into searchable chunks (embeddings in a vector store).
2. For each coding task you **retrieve** the closest chunks (a User handler, an Auth JWT claim, a Mongo repository, a paragraph from this README).
3. You send **those chunks + a narrow question** to the LLM.
4. The model writes Driver code that follows *your* patterns instead of a generic tutorial.

That is the whole point of using RAG as the foundation: **ground generation in this taxi system**, not in the internet average.

You still own the result. RAG is a junior pair-programmer with your repo open. It is not a merge button.

```text
  ┌──────────────┐     embed      ┌─────────────┐
  │ Spec README  │───────────────►│             │
  │ User Service │                │ Vector store│
  │ Auth Service │───────────────►│ (Qdrant)    │
  └──────────────┘                └──────┬──────┘
                                         │ top-k chunks
  “toggle driver status, same            ▼
   layering as User profile”  →  LLM  →  draft Go
                                         │
                                         ▼
                                   you review, test, commit
```

---

## What you must already have (from User + Auth)

Do not start this milestone until all of this is true:

- [ ] User Service: register / login-support / profile, 3-layer layout (handler → service → repository), MongoDB, Gin
- [ ] Auth Service: access + refresh JWT, Redis, `/refresh`, logout
- [ ] A user with role `Driver` can be created through User Service (even if Driver Service is still a stub)
- [ ] Shared conventions exist in the repo and you can point to them: error JSON, config via env, Docker Compose, `golangci-lint`, test layout
- [ ] You can explain every important file in User + Auth. If you cannot, RAG will only help you copy code you do not understand

Bring these **inputs** into the RAG corpus (see below). They are the source of truth, not a new blog post about “how to write a microservice”.

---

## What you need to set up

Only if you take this optional track. Local-first. Cloud APIs are optional. **No API keys in git.**

### 1. LLM (the generator)

Pick **one**:

| Option | When to use | What to run |
| --- | --- | --- |
| **Ollama (recommended)** | default; works offline; no card | install Ollama, pull a **coding** model |
| OpenAI / Anthropic / etc. | only if you already have a key and the mentor agrees | env var in a **gitignored** `.env` |

Ollama:

```bash
# install: https://ollama.com/download
ollama pull qwen2.5-coder:7b    # generator (14b if your RAM allows)
ollama pull nomic-embed-text    # embeddings for RAG
ollama serve                    # if it is not already running
```

Sanity check: `ollama run qwen2.5-coder:7b "write a Go function that returns an error, nothing else"`.

### 2. Embeddings (the retriever’s language)

Same space for index and query. Do not embed with model A and query with model B.

- Local: `nomic-embed-text` via Ollama (768-d)
- Cloud: OpenAI `text-embedding-3-small` only if the LLM is also cloud

### 3. Vector store (where chunks live)

| Option | When to use | What to run |
| --- | --- | --- |
| **Qdrant** (recommended) | you already live in Docker Compose | `docker compose` service on `6333` |
| [chromem-go](https://github.com/philippgille/chromem-go) | zero extra containers | embedded in a small Go indexer |

Example Compose service (add next to Mongo/Redis, do not publish it outside localhost):

```yaml
qdrant:
  image: qdrant/qdrant:latest
  ports:
    - "127.0.0.1:6333:6333"
  volumes:
    - qdrant_data:/qdrant/storage
```

Sanity check: open `http://127.0.0.1:6333/dashboard` or `curl -s http://127.0.0.1:6333/readyz`.

### 4. Indexer + ask loop (you write this)

A small tool in the monorepo, for example `tools/rag/` (Go is preferred). It must:

1. **Chunk** the corpus (see next section): ~400–800 tokens, overlap ~50–100, **keep file path + start line** in metadata
2. **Embed** each chunk and upsert into Qdrant (collection e.g. `inno-taxi`)
3. **Ask**: embed the question → search top-k (k = 6–12) → build a prompt → call Ollama/OpenAI → print answer **and** the retrieved file paths

Re-index after User/Auth changes. Stale index = confident wrong code.

You may use [langchaingo](https://github.com/tmc/langchaingo) or call Ollama HTTP (`/api/embed`, `/api/chat`) and Qdrant REST yourself. Either is fine. A pile of untitled ChatGPT tabs is not.

### 5. Driver Service runtime (same as the spec)

RAG does not replace Mongo or Gin. The **running** Driver Service still needs:

- MongoDB (same family as User Service; separate database/collection, not the User DB)
- Gin, mongo-driver, golang-migrate, testify
- Compose entry for `driver-service`
- The User Service must be able to **call** Driver Service on driver registration (gRPC preferred; HTTP is acceptable for this milestone if gRPC is not ready yet)

Trip accept/decline that needs a live Order Service: **stub the port** (interface + fake). Do not wait for Order Service.

---

## What to index (the corpus)

Index **only** what should constrain the model. Garbage in, taxi-app-shaped garbage out.

| Include | Why |
| --- | --- |
| This folder: `README.md` (Driver / User / Auth / NFR / auth flow) | product rules |
| This file | milestone rules |
| `user-service/` (handlers, services, repos, models, errors, tests) | copy **structure**, not business fields blindly |
| `auth-service/` (JWT claims, token payload, middleware) | Driver is a role on the same tokens |
| `docker-compose.yml`, Makefiles, `.golangci.yml` | how we run and lint |
| proto / OpenAPI **already in the repo** | contracts |

| Exclude | Why |
| --- | --- |
| `vendor/`, `node_modules/`, binaries, `*.pb.go` if huge | noise |
| `.env`, keys, dumps | secrets |
| `go.sum`, lockfiles | zero design signal |
| random RAG tutorials | the model already knows those; they **fight** your house style |

Chunk by file, preferably by function/type. A 2 000-line `handler.go` as one blob retrieves poorly.

---

## How a single task is done (work the loop, do not one-shot the service)

Do **one vertical slice per query**, then test, then the next slice.

Suggested order:

1. Skeleton: `cmd/`, `internal/{handler,service,repository}`, config, health, Compose
2. Driver document in Mongo (fields that are **driver-only**; identity stays in User Service)
3. Internal **register** endpoint that User Service calls when role is `Driver`
4. Profile get/update (JWT, role `Driver` or `Admin` — Gateway will enforce later; still check claims if the service is called directly)
5. Status: `offline` \| `available` \| `on-trip` with allowed transitions
6. Ports for Order Service: `AcceptTrip`, `DeclineTrip`, `StartTrip`, `CompleteTrip` — interface + in-memory fake + tests
7. OpenAPI or proto, README, tests

### Prompt shape (use this every time)

```text
You are implementing Driver Service in the inno-taxi monorepo.

RULES
- Follow retrieved code: layering, package names, error JSON, env config, logging.
- Stack: Go, Gin, mongo-driver, MongoDB. No GORM, no Postgres, no Fiber.
- Identity lives in User Service. Driver Service stores driver-only data.
- If retrieved context is missing a fact, say WHAT IS MISSING instead of inventing a stack.

TASK
<one slice, e.g. PATCH status>

RETRIEVED CONTEXT
<paste the chunks your tool returned, with file paths>
```

Bad prompt: “Write a complete driver microservice.”  
Good prompt: “Add `PATCH /drivers/me/status` with body `{"status":"available"}`. Mirror User Service profile handler: bind JSON, call service, map domain errors to the same HTTP error shape. Default status on register is `offline`. Forbidden: `offline` → `on-trip`.”

### Worked example — why retrieval changes the output

**Question you type into `tools/rag`:**

> Implement register-from-User-Service: when User Service creates a user with role Driver, it must call Driver Service to insert a driver profile. Default status `offline`. Follow our handler → service → repository split and our Mongo error handling.

**What a useful retrieval looks like** (illustrative paths — yours will differ):

| Score | Chunk | Why it was retrieved |
| --- | --- | --- |
| high | `user-service/internal/handler/register.go` | same “create related aggregate” moment |
| high | `user-service/internal/repository/user_mongo.go` | how you insert and wrap duplicate-key |
| high | README “Driver Service / Registration” | product rule: User Service invokes Driver |
| medium | `auth-service/.../claims.go` | `user_id` + roles in the JWT |
| medium | User Docker Compose Mongo URI | how services find the DB |

**What the LLM should produce:** an internal handler (not public signup), a `Driver` document keyed by `user_id`, repository insert, tests for duplicate `user_id`.

**What a model *without* RAG often produces (reject this):** public `POST /drivers/register` with email+password, bcrypt, Postgres, a new JWT issuer. That duplicates Auth/User and breaks the spec.

**What you do after the draft:**

1. Read every line. Delete anything that duplicates User/Auth (login, password hash, issuing JWT).
2. Wire User Service to call Driver Service **inside the registration path** (or outbox/event if you already have Kafka; sync call is enough for this milestone). If Driver Service is down, registration of a Driver must **fail** (no half-created drivers). Document that choice.
3. `go test ./...` on Driver Service; hit the happy path with curl/Postman.
4. Save the query + retrieved paths in `driver-service/rag/queries.md` (see expected result).

If top-k chunks are irrelevant (e.g. wallet comments, unrelated tests), **fix the question or the chunking**, do not “just generate anyway”.

---

## Guardrails

- One slice per LLM call. A 20-file dump is how you get inconsistent layers.
- If the model invents a library that is not in User/Auth/`go.mod`, do not add it “because the LLM said so”.
- Generated tests that only assert `assert.True(t, true)` are not tests.
- You must be able to explain status transitions and why identity is not stored twice.
- Cursor / Copilot may help **after** the corpus is indexed. They do not replace `tools/rag` **on this optional track** — the mentor will ask how retrieval worked.

---

## Expected result

A mentor can clone the monorepo, follow Driver README, and get this without guessing.

### Service (product)

- [ ] Driver Service runs via Docker Compose next to User + Auth + Mongo + Redis (+ Qdrant for the indexer, if you used it)
- [ ] Clean architecture: handler / service / repository; no DB calls from handlers
- [ ] Stack matches the spec: Gin, mongo-driver, MongoDB, testify
- [ ] **Register:** User Service creates a user with role `Driver` → Driver profile exists with `user_id`, status `offline`. Repeat register for the same `user_id` is a documented conflict, not a second row
- [ ] **Profile:** driver can get/update **driver-only** fields (e.g. taxi type Economy/Comfort/Business, license plate). Name/email/phone stay in User Service
- [ ] **Status:** `offline` / `available` / `on-trip`; illegal transitions return 409 (or your existing domain-error mapping — **same as User Service**, not a new envelope)
- [ ] Order-facing methods exist as a **Go interface** + fake; HTTP/gRPC handlers may be stubbed with `501` or a clear TODO **only** for live assignment from Order Service
- [ ] No password hashing and no JWT issuance in Driver Service
- [ ] README: how to run, how to test, env vars, how User Service finds Driver Service

### RAG (process — only if you opted in)

If you skipped RAG, ignore this subsection. The product checklist above is still required.

If you opted in, commit a folder, e.g. `driver-service/rag/`:

| File | Contents |
| --- | --- |
| `sources.md` | list of indexed paths, embed model name, vector DB, chunk size |
| `queries.md` | each slice: question, **file paths retrieved**, 3–6 line note: kept / edited / rejected |
| `NOTES.md` | one short paragraph: a case where RAG was wrong and you fixed it by hand |

Without `rag/`, you simply did not take the optional track. That is fine. If you *did* opt in, `rag/` is how the mentor sees that retrieval actually happened — a Driver Service that “looks generated” with no corpus notes does not count as this track.

### Quality bar (same as User + Auth)

- `gofmt` / `golangci-lint` clean
- unit tests for status transitions and duplicate `user_id`
- at least one integration test with Mongo (dockertest or Compose)
- secrets only in `.env` (gitignored)

---

## Where to read (do this before writing `tools/rag`)

Read **in this order**. Official docs over random Medium posts.

### Ideas (short)

1. [Retrieval-augmented generation — Wikipedia](https://en.wikipedia.org/wiki/Retrieval-augmented_generation) — vocabulary: retrieve, then generate  
2. [Anthropic: Contextual retrieval](https://www.anthropic.com/news/contextual-retrieval) — why naive chunking fails; you do not need their product  
3. [Pinecone: What is RAG](https://www.pinecone.io/learn/retrieval-augmented-generation/) — embeddings + top-k, diagrams  

Optional depth: [OpenAI — RAG](https://help.openai.com/en/articles/8868588-retrieval-augmented-generation-rag-and-semantic-search) if you use their API.

### Things you actually run

4. [Ollama](https://github.com/ollama/ollama) — install, `pull`, `/api/embed`, `/api/chat`  
5. [Qdrant quickstart](https://qdrant.tech/documentation/quickstart/) — collections, upsert, search  
6. [nomic-embed-text](https://ollama.com/library/nomic-embed-text) — local embeddings  

### Go

7. [langchaingo](https://github.com/tmc/langchaingo) — optional; Ollama + vector stores  
8. [chromem-go](https://github.com/philippgille/chromem-go) — optional embedded store  
9. Your **User Service** — this is the main style guide. Read it like a library.

### Not required

LangGraph, agents with 12 tools, fine-tuning, “build a chatbot UI”. The driver app does not need an LLM **at runtime**. If you opt in, RAG is how **you** write the service.

---

## Suggested time (extra, on top of writing Driver Service)

Skip this table if you are not doing RAG.

| Block | Time |
| --- | --- |
| Read the RAG links + sketch chunk/ask | 0.5–1 day |
| `tools/rag` index + one demo query that retrieves `register.go` | 1–2 days |
| Driver slices 1–5 with the loop | 4–8 days |
| Tests, Compose, `rag/` writeup, PR | 1–2 days |

If retrieval never returns User Service files, stop generating and fix the indexer (or drop the optional track and finish Driver Service by hand).
