# Intern brief — RAG that grades drivers from passenger comments

**You will build a small RAG pipeline.** After a trip, passengers leave **comments**. Your job is to retrieve those comments, let an LLM read them **together with a fixed rubric**, and set each driver a grade: **`low` | `medium` | `perfect`**.

This is **not** “ask ChatGPT to write Driver Service.”  
This is **not** a chatbot UI.  
The LLM runs **when you grade a driver**, not on every HTTP request of the taxi app.

---

## What you are building (one sentence)

Index passenger comments → for a given `driver_id`, retrieve that driver’s comments → send comments + rubric to the LLM → write `grade` on the driver document.

```text
  ┌─────────────────────┐   embed    ┌─────────────┐
  │ Comments (corpus)   │───────────►│ Vector store│
  │ Grade rubric (this  │            │ (Qdrant)    │
  │ file)               │───────────►│             │
  └─────────────────────┘            └──────┬──────┘
                                            │ top-k comments
  “Grade driver <id>”  ──embed query──► search
                                            ▼
                                   LLM + rubric
                                            ▼
                              grade: low | medium | perfect
                                            ▼
                              save on Driver (Mongo)
```

---

## Why RAG (and not “paste every comment into ChatGPT”)

**Retrieval-Augmented Generation** means: **find the right text first**, then generate.

| Without retrieval | With RAG (what we want) |
| --- | --- |
| You dump 2 000 comments into the prompt, hit token limits, or the model ignores old trips | You retrieve the **most relevant** comments for *this* driver (and the rubric) |
| The model invents a 1–5 star scale or writes an essay | The model must answer with **only** `low`, `medium`, or `perfect` |
| Grades drift every time you ask | Same rubric is always in the prompt; comments come from **your** store, not the internet |

You still own the result. If the retrieved comments are wrong, the grade is wrong — **fix indexing**, do not “just generate anyway.”

---

## Grade rubric (source of truth — index this)

The LLM must **not** invent extra labels (`excellent`, `3 stars`, `toxic`). Only these three:

| Grade | When to use |
| --- | --- |
| **`low`** | Comments are mostly negative: rude, unsafe driving, no-show, insults, “never again”, dirty car, cheating on route/price. A few polite comments do **not** cancel a pattern of harm. |
| **`medium`** | Mixed or ordinary: on-time but quiet, small complaints (AC, music) plus some thanks, not enough evidence for `perfect` or `low`. **Default when comments are few or vague** (“ok”, “fine”). |
| **`perfect`** | Comments are consistently positive: polite, safe, on time, clean car, helpful with bags, “best driver”. No serious safety/rudeness complaints in the retrieved set. |

**Tie-breakers (put these in the prompt every time):**

1. **Safety and rudeness beat praise.** One “he almost crashed” + ten “nice guy” → not `perfect`; usually `low`.
2. **Need a pattern for `low`/`perfect`.** One short “ok” comment → `medium`.
3. **If there are no comments** → do **not** call the LLM. Return `ungraded` (or leave `grade` empty) and a clear error: `no_comments`.
4. **Output JSON only** (see prompt below). No markdown, no extra keys.

Numeric trip stars (1–5), if you have them, are **optional extra context**. The grade is decided from **comment text**, not by averaging stars alone.

---

## What you must already have

Do not start until this is true:

- [ ] You know what a **driver** is in this project (Driver Service stores driver-only data; identity lives in User Service)
- [ ] You understand that **comments live with trips/orders** (Order Service: completed trips can have a rating **and** a comment). If Order Service is not ready, you **seed mock comments** (see below) — do not block on Kafka
- [ ] Docker Compose (or Docker Desktop) works on your machine
- [ ] You can explain RAG in your own words: embed → search → prompt → generate

---

## What you need to set up (local-first)

**No API keys in git.** Cloud LLMs are optional and only if a mentor agrees.

### 1. LLM (the grader)

Pick **one**:

| Option | When to use | What to run |
| --- | --- | --- |
| **Ollama (recommended)** | default; works offline | install Ollama, pull a small chat model |
| OpenAI / Anthropic | only with a key + mentor OK | env var in a **gitignored** `.env` |

```bash
# install: https://ollama.com/download
ollama pull llama3.2          # grader (or qwen2.5:7b if you have RAM)
ollama pull nomic-embed-text  # embeddings — required for RAG
ollama serve                  # if it is not already running
```

Sanity check:

```bash
ollama run llama3.2 "Reply with only the word medium"
```

### 2. Embeddings (same model for index and query)

Do not embed comments with model A and search with model B.

- Local: `nomic-embed-text` via Ollama (768 dimensions)
- Cloud: `text-embedding-3-small` only if the LLM is also cloud

### 3. Vector store

| Option | When to use | What to run |
| --- | --- | --- |
| **Qdrant (recommended)** | you already use Docker Compose | service on `6333`, bind to localhost |
| [chromem-go](https://github.com/philippgille/chromem-go) | no extra container | embedded in a small Go tool |

```yaml
qdrant:
  image: qdrant/qdrant:latest
  ports:
    - "127.0.0.1:6333:6333"
  volumes:
    - qdrant_data:/qdrant/storage
```

Sanity check: `http://127.0.0.1:6333/dashboard` or `curl -s http://127.0.0.1:6333/readyz`.

### 4. Tool you write: `tools/driver-grade-rag/` (Go preferred)

A small CLI or HTTP helper in the monorepo. It must:

1. **Chunk / index comments** — one comment = one chunk (do not glue 50 trips into one blob). Metadata **must** include: `driver_id`, `comment_id` (or trip/order id), `created_at`, optional `star_rating`
2. **Embed** each comment and upsert into Qdrant (collection e.g. `driver-comments`)
3. **Grade:** take `driver_id` → embed a query like `passenger comments about driver {id}` **or** filter Qdrant by `driver_id` (filter is better) → take top-k (k = 8–20) → build prompt with rubric + comments → call Ollama → parse JSON → print grade **and** the comment ids used
4. **Persist:** write `grade` + `graded_at` + `comment_ids_used` onto the driver (Driver Service Mongo, or a local JSON file if Driver Service is not ready — document which)

You may use [langchaingo](https://github.com/tmc/langchaingo) or call Ollama HTTP (`/api/embed`, `/api/chat`) and Qdrant REST yourself. Untitled ChatGPT tabs are not a submission.

Re-index when comments change. Stale index = confident wrong grade.

---

## What to index (the corpus)

Index **only** what should affect the grade.

| Include | Why |
| --- | --- |
| Passenger **comments** (text) keyed by `driver_id` | this is the evidence |
| This rubric section | the model must not invent a scale |
| Optional: 5–10 **labelled examples** you wrote (comment → expected grade) | few-shot, keeps the model consistent |

| Exclude | Why |
| --- | --- |
| Driver passwords, JWTs, `.env` | secrets |
| Full User/Auth source code | noise; this task is not “generate a microservice” |
| Random RAG blog posts | they fight your rubric |
| Other drivers’ comments **without** a `driver_id` filter | contaminates the grade |

If Order Service is missing, commit **seed data**, e.g. `tools/driver-grade-rag/testdata/comments.json`:

```json
[
  {
    "comment_id": "c1",
    "driver_id": "drv_01",
    "text": "Drove like a maniac, I was scared.",
    "star_rating": 1,
    "created_at": "2026-08-01T10:00:00Z"
  },
  {
    "comment_id": "c2",
    "driver_id": "drv_01",
    "text": "Never again. Rude and on the phone.",
    "star_rating": 1,
    "created_at": "2026-08-02T10:00:00Z"
  },
  {
    "comment_id": "c3",
    "driver_id": "drv_02",
    "text": "On time, quiet, car was clean. Thank you!",
    "star_rating": 5,
    "created_at": "2026-08-03T10:00:00Z"
  }
]
```

Expected: `drv_01` → `low`, `drv_02` → `perfect` (if that is their only comment, `perfect` is allowed; if you only have one vague “thanks”, `medium` is safer — **follow the rubric**, and write your choice in `NOTES.md`).

---

## How to implement (do not one-shot the whole pipeline)

Do **one slice**, test it, then the next.

1. **Seed comments** + Qdrant collection (payload schema with `driver_id`)
2. **Index:** read JSON (later: Mongo/Order) → embed → upsert. Log how many vectors
3. **Retrieve-only command:** `grade-rag retrieve --driver drv_01` prints comments, **no LLM**. You must see only that driver’s texts
4. **Grade command:** retrieve → prompt → parse JSON → print `low|medium|perfect`
5. **Save grade** on the driver document (field `grade`)
6. **README** for the tool: how to run, env vars, how to re-index

### Prompt shape (use this every time)

```text
You grade taxi drivers from passenger comments.

RULES
- Use ONLY the retrieved comments and the rubric.
- Reply with JSON only, no markdown:
  {"grade":"low"|"medium"|"perfect","reason":"<one short sentence>"}
- If comments are missing or empty, you should not be called. If you still are, return grade "medium" and reason "insufficient evidence".
- Do not mention other drivers. Do not invent comments.

RUBRIC
<paste the grade table + tie-breakers>

RETRIEVED COMMENTS FOR DRIVER <driver_id>
<each comment: id, date, optional stars, text>
```

Bad query: “Is this driver good?”  
Good query: “Grade `drv_01` using only retrieved comments; output the JSON schema above.”

### Worked example

**You run:**

```bash
go run ./tools/driver-grade-rag grade --driver drv_01
```

**Useful retrieval:**

| Payload | Why it is useful |
| --- | --- |
| `c1` “Drove like a maniac…” | safety → pulls toward `low` |
| `c2` “Rude and on the phone” | pattern, not a one-off |

**What the LLM should return:**

```json
{"grade":"low","reason":"Multiple comments report unsafe driving and rudeness."}
```

**What a model without RAG often returns (reject this):** `"4.2/5"`, `"needs coaching"`, or a grade based on **another** driver’s comments because you forgot the `driver_id` filter.

**What you do after:**

1. Read the retrieved comments. If they are the wrong driver, **stop** — fix filter/metadata
2. Check JSON parse; invalid JSON = retry once, then fail the command (do not guess `perfect`)
3. Save `grade` on the driver
4. Log in `tools/driver-grade-rag/queries.md`: driver id, comment ids retrieved, grade, kept / edited

---

## Guardrails

- Filter by `driver_id` in Qdrant. Semantic search alone can mix drivers with similar wording
- Do not send PII you do not need (full passenger names, phones). Comment text + ids are enough
- Do not fine-tune a model. Do not build LangGraph “agents”
- Generated “tests” that only `assert.Equal(t, true, true)` do not count
- You must be able to explain: embeddings, top-k, why `medium` is the default for weak evidence

---

## Expected result

A mentor clones the repo, follows **your** tool README, and gets a grade without guessing.

### Product

- [ ] Comments can be indexed (seed JSON is enough for the milestone)
- [ ] `retrieve --driver <id>` returns **only** that driver’s comments
- [ ] `grade --driver <id>` prints valid JSON with `grade` ∈ {`low`,`medium`,`perfect`} and a one-line `reason`
- [ ] Grade is stored on the driver (`grade`, `graded_at`, list of `comment_ids_used`)
- [ ] No comments → no LLM call; explicit `no_comments` / ungraded
- [ ] Seed fixtures: at least one driver that grades `low` and one that grades `perfect` or `medium` (document expected grades)

### Process (so we see it was actually RAG)

Commit `tools/driver-grade-rag/` (name may differ) including:

| File | Contents |
| --- | --- |
| `README.md` | run, re-index, env, Ollama model names |
| `sources.md` | what you indexed, embed model, vector DB, chunk = one comment |
| `queries.md` | a few `driver_id`s: retrieved comment ids, model grade, whether you overrode it |
| `NOTES.md` | one case where retrieval or the model was wrong and you fixed it |
| `testdata/` | seed comments |

### Quality bar

- Deterministic CLI (exit non-zero on Qdrant/Ollama down)
- At least one test: given fixture comments for `drv_01`, retrieval returns those ids
- Secrets only in `.env` (gitignored)

---

## Where to read (before writing code)

Read **in this order**. Official docs over random Medium posts.

1. [Retrieval-augmented generation — Wikipedia](https://en.wikipedia.org/wiki/Retrieval-augmented_generation) — retrieve, then generate
2. [Pinecone: What is RAG](https://www.pinecone.io/learn/retrieval-augmented-generation/) — embeddings + top-k
3. [Ollama](https://github.com/ollama/ollama) — `pull`, `/api/embed`, `/api/chat`
4. [Qdrant quickstart](https://qdrant.tech/documentation/quickstart/) — collections, payload **filters**, search
5. [nomic-embed-text](https://ollama.com/library/nomic-embed-text)

Optional: [langchaingo](https://github.com/tmc/langchaingo), [chromem-go](https://github.com/philippgille/chromem-go).

**Not required:** LangGraph, fine-tuning, a chat UI, generating Driver Service from User/Auth.

---

## Suggested time

| Block | Time |
| --- | --- |
| Read this brief + RAG links; sketch payload schema | 0.5 day |
| Qdrant + Ollama + index seed comments | 1 day |
| Retrieve-by-`driver_id` (no LLM) + tests | 0.5–1 day |
| Grade prompt + JSON parse + save `grade` | 1–2 days |
| Writeup (`queries.md`, `NOTES.md`), PR | 0.5 day |

If retrieval returns another driver’s comments, **stop grading** and fix metadata/filters.
