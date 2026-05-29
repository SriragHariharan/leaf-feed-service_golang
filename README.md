# Feed Service (Go)

## Service Name & Overview

The Feed Service serves cursor-paginated home feeds and user timelines. It maintains feed state in CockroachDB/PostgreSQL and consumes user and interaction events from Kafka to keep feeds up to date with likes, comments, and profile changes.

**Tech Stack**

- **Language:** Go 1.26.2
- **Framework:** Gorilla Mux
- **ORM:** GORM (PostgreSQL / CockroachDB driver)
- **Key libraries:** segmentio/kafka-go, golang-jwt/jwt, godotenv

## Architecture & Dependencies

### Internal Dependencies

| Dependency | Purpose |
|---|---|
| **CockroachDB / PostgreSQL** | Feed, post, and user tables via GORM |
| **Kafka** | Consumes user sync and interaction events |

### Event Contracts

See [`../KAFKA_TOPICS.txt`](../KAFKA_TOPICS.txt) for the platform topic list.

| Direction | Topic | Consumer group | Events |
|---|---|---|---|
| **Produces** | — | — | None |
| **Consumes** | `user.events` | `feed-service-go-user-events` | `{ userID, username, profilePicture }` |
| **Consumes** | `interaction.events` | `feed-service-go-interaction-events` | `post.liked`, `post.unliked`, `post.commented`, `post.uncommented` |

### External APIs

None.

## Environment Variables

```bash
# --- Server ---
PORT=2004

# --- Kafka (local — plaintext Docker) ---
KAFKA_MODE=local
KAFKA_BROKERS=localhost:9092

# --- Kafka (Aiven — uncomment and set KAFKA_MODE=aiven) ---
# KAFKA_MODE=aiven
# KAFKA_BROKERS=your-service.a.aivencloud.com:12345
# KAFKA_SASL_USERNAME=your-aiven-username
# KAFKA_SASL_PASSWORD=your-aiven-password
# KAFKA_SASL_MECHANISM=scram-sha-256
# KAFKA_SSL_CA_PATH=./ca.pem
# KAFKA_SSL_CA=

# --- Auth / JWT ---
ACCESS_TOKEN_SECRET=your-access-token-secret

# --- Database (CockroachDB / PostgreSQL) ---
# COCKROACHDB_* are documentation helpers; runtime uses DATABASE_URL only
COCKROACHDB_PASSWORD=your-cockroach-password
COCKROACHDB_USER=your-cockroach-user
DATABASE_URL=postgresql://user:password@host:26257/defaultdb?sslmode=verify-full

# --- App (documented; not read by Go runtime) ---
SERVICE_NAME=feed-service-go
NODE_ENV=dev
LOG_LEVEL=info
APP_VERSION=1.0.0

# --- Optional: connection pool tuning ---
# DB_MAX_OPEN_CONNS=25
# DB_MAX_IDLE_CONNS=10
# DB_CONN_MAX_LIFETIME=30m
# DB_CONN_MAX_IDLE_TIME=5m
# GORM_LOG_LEVEL=warn

# --- Optional: seed CLI only ---
# SEED_TRUNCATE=1
# SEED_RNG_SEED=42
```

> **Cross-service note:** `ACCESS_TOKEN_SECRET` must match the value configured in user-service.

## Getting Started

### Prerequisites

- **Go** 1.26.2+
- **CockroachDB** or **PostgreSQL** instance (local or cloud)
- **Kafka** from parent docker-compose
- **Air** (optional, for live reload): `go install github.com/air-verse/air@latest`

### Local Infrastructure

```bash
# From the parent repo root (d:\main PROJECTS\leaf\)
docker compose up -d kafka
```

Set up a CockroachDB or PostgreSQL database and configure `DATABASE_URL` in `.env`.

### Install & Run

```bash
cd feed-service-go
cp .env.example .env
# Edit .env with your DATABASE_URL and ACCESS_TOKEN_SECRET
go mod download
go run ./cmd/server
```

Verify: `curl http://localhost:2004/test` → `Feed service running!`

**Live reload (optional):**

```bash
air
```

**Seed test data (optional):**

```bash
go run ./cmd/seed
# Destructive re-seed:
SEED_TRUNCATE=1 go run ./cmd/seed
```

## Available Scripts

This service has no `package.json`. Use Go commands directly:

| Command | Description |
|---|---|
| `go run ./cmd/server` | Start the HTTP server and Kafka consumers |
| `go run ./cmd/seed` | Seed the database with test data |
| `air` | Live reload during development (via `.air.toml`) |
| `go build -o feed-service ./cmd/server` | Build a production binary |
| `go test ./...` | Run all tests |

## API / Event Interface

Gateway prefixes: `/api/v1/feed`, `/api/v1/timeline`

### HTTP Routes

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/test` | No | Health check (not proxied through gateway) |
| `GET` | `/api/v1/feed` | Yes | Cursor-paginated home feed (`?cursor=`) |
| `GET` | `/api/v1/timeline/{user_id}` | Yes | User timeline; `{user_id}` can be `self` |

All authenticated routes require a Bearer JWT in the `Authorization` header.

### Event-Driven Behavior

This service does not expose Kafka endpoints. It runs background consumers that update feed state in response to:

- **`user.events`** — upserts user profile data used in feed rendering
- **`interaction.events`** — updates interaction flags (liked, commented) on feed entries
