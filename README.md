# Full-Stack Calculator Application

A full-stack calculator application featuring a Go REST API microservice backend, a modern React 19 and TypeScript frontend, and Docker Compose orchestration with an Nginx reverse proxy.

---

## Table of Contents
- [Architecture & Tech Stack](#architecture--tech-stack)
- [Features](#features)
- [Project Layout](#project-layout)
- [Quick Start](#quick-start)
  - [Option 1: Docker Compose (Recommended)](#option-1-docker-compose-recommended)
  - [Option 2: Local Development](#option-2-local-development)
- [REST API Specification](#rest-api-specification)
  - [Endpoints](#endpoints)
  - [API Examples (cURL)](#api-examples-curl)
  - [Error Handling & Status Codes](#error-handling--status-codes)
- [Testing & Coverage](#testing--coverage)
  - [Backend Tests](#backend-tests)
  - [Frontend Tests](#frontend-tests)
- [Design Decisions & Assumptions](#design-decisions--assumptions)

---

## Architecture & Tech Stack

```
┌────────────────────────────────────────────────────────┐
│                   Client Browser                       │
│        (React 19 + TypeScript + Vitest + CSS)          │
└──────────────┬─────────────────────────┬───────────────┘
               │ (Docker Mode via /api)  │ (Local Dev Mode)
               ▼                         │
┌──────────────────────────────┐         │
│     Nginx Reverse Proxy      │         │
│  (Port 3000 / Alpine Linux)  │         │
└──────────────┬───────────────┘         │
               │                         ▼
               └──────────────► ┌────────────────────────────────┐
                                │          Go REST API           │
                                │   (Port 8080 / Go 1.24 Alpine) │
                                └────────────────────────────────┘
```

- **Backend**: Go 1.24 using standard library `net/http` (Go 1.22+ routing patterns), zero external runtime dependencies, structured logging with `log/slog`, strict JSON validation, and panic recovery middleware.
- **Frontend**: React 19, TypeScript, Vite, CSS Grid layout with custom hardware styling, responsive mobile layout, accessible ARIA live regions, and `localStorage` history tape.
- **Containerization**: Multi-stage Docker builds (`golang:1.24-alpine` -> `alpine:3.21` running as non-root user; `node:22-alpine` -> `nginx:1.27-alpine`), unified through `docker-compose.yml` with health check dependencies.

---

## Features

- **Standard & Advanced Arithmetic**: Addition, Subtraction, Multiplication, Division, Exponentiation (`^`), Square Root (`√`), and Percentage (`%`).
- **Precision Floating-Point Normalization**: Prevents IEEE 754 precision artifacts (e.g., `0.1 + 0.2 = 0.3` instead of `0.30000000000000004`) and normalizes negative zero (`-0` to `0`).
- **History Tape**: Persists calculation history to `localStorage` with timestamping and one-click recall of previous calculations.
- **Keyboard Support**: Full physical keyboard input (`0-9`, `.`, `+`, `-`, `*`, `/`, `^`, `%`, `Enter`/`=`, `Backspace`/`Delete`, `Escape`/`C`).
- **Responsive UI**: Retro hardware aesthetic with dual-line display, tabular numbers to eliminate font jitter, and fluid layouts for both desktop and mobile screens.
- **Robust Error Handling**: Server-side validation for division by zero, negative square roots, malformed payloads, payload size caps (1MB), and unknown JSON keys.

---

## Project Layout

```
.
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Server bootstrap, graceful shutdown, CORS
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── service.go           # Arithmetic domain logic and precision rounding
│   │   │   └── service_test.go      # Table-driven unit tests
│   │   ├── handler/
│   │   │   ├── handler.go           # HTTP decoding, validation, error formatting
│   │   │   └── handler_test.go      # Handler tests with net/http/httptest
│   │   └── server/
│   │       ├── router.go            # ServeMux routing, recovery, logging, CORS
│   │       └── router_test.go       # End-to-end router and middleware tests
│   ├── Dockerfile                   # Multi-stage minimal Alpine runner
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   ├── client.ts            # Typed HTTP client for calculator endpoints
│   │   │   └── types.ts             # DTO interfaces and operator types
│   │   ├── components/
│   │   │   ├── Calculator.tsx       # Main chassis, keyboard listener, layout
│   │   │   ├── Calculator.test.tsx  # Component integration tests
│   │   │   ├── Display.tsx          # Dual-line LCD display with status indicators
│   │   │   ├── Keypad.tsx           # Grid layout keypad with ARIA labels
│   │   │   └── History.tsx          # Collapsible history tape drawer
│   │   ├── hooks/
│   │   │   ├── useCalculator.ts     # Calculator state machine and history persistence
│   │   │   └── useCalculator.test.ts# Hook unit tests
│   │   ├── utils/
│   │   │   └── format.ts            # String formatting, digit buffers, sign toggle
│   │   ├── setupTests.ts            # Storage mocks and jest-dom matchers
│   │   └── index.css                # Component styling, grid layout, theme tokens
│   ├── nginx.conf.template          # Reverse proxy template for /api and /health
│   ├── Dockerfile                   # Multi-stage Vite build + Nginx Alpine
│   └── package.json
├── docker-compose.yml               # Multi-container orchestration & networking
├── PROMPTS.md                       # AI collaboration log
└── README.md
```

---

## Quick Start

### Option 1: Docker Compose (Recommended)

Run the entire stack with a single command:

```bash
docker compose up --build
```

- **Frontend Application**: [http://localhost:3000](http://localhost:3000)
- **Backend API Direct**: [http://localhost:8080](http://localhost:8080)
- **Backend Health Check**: [http://localhost:8080/health](http://localhost:8080/health) (or proxied via [http://localhost:3000/health](http://localhost:3000/health))

To shut down the containers:
```bash
docker compose down
```

---

### Option 2: Local Development

#### Prerequisites
- Go 1.22+
- Node.js 20+ and npm

#### 1. Start the Backend

```bash
cd backend
go run ./cmd/api
```
The backend starts on `http://localhost:8080` (default CORS allows `http://localhost:5173`).

#### 2. Start the Frontend

In a separate terminal:

```bash
cd frontend
npm install
npm run dev
```
The frontend dev server starts on `http://localhost:5173`.

---

## REST API Specification

### Endpoints

| Method | Endpoint | Description | Request Body |
|---|---|---|---|
| `POST` | `/api/v1/add` | Addition ($a + b$) | `{"a": number, "b": number}` |
| `POST` | `/api/v1/subtract` | Subtraction ($a - b$) | `{"a": number, "b": number}` |
| `POST` | `/api/v1/multiply` | Multiplication ($a \times b$) | `{"a": number, "b": number}` |
| `POST` | `/api/v1/divide` | Division ($a / b$) | `{"a": number, "b": number}` |
| `POST` | `/api/v1/power` | Exponentiation ($a^b$) | `{"a": number, "b": number}` |
| `POST` | `/api/v1/sqrt` | Square Root ($\sqrt{a}$) | `{"a": number}` |
| `POST` | `/api/v1/percentage`| Percentage ($a / 100$) | `{"a": number}` |
| `GET` | `/health` | Service health status | _None_ |

---

### API Examples (cURL)

#### 1. Addition (Binary Operation)
```bash
curl -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" \
  -d '{"a": 12.5, "b": 7.5}'
```
**Response (200 OK):**
```json
{
  "result": 20
}
```

#### 2. Square Root (Unary Operation)
```bash
curl -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" \
  -d '{"a": 144}'
```
**Response (200 OK):**
```json
{
  "result": 12
}
```

#### 3. Power (Exponentiation)
```bash
curl -X POST http://localhost:8080/api/v1/power \
  -H "Content-Type: application/json" \
  -d '{"a": 2, "b": 8}'
```
**Response (200 OK):**
```json
{
  "result": 256
}
```

#### 4. Division by Zero (Validation Error)
```bash
curl -X POST http://localhost:8080/api/v1/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 42, "b": 0}'
```
**Response (422 Unprocessable Entity):**
```json
{
  "error": "division by zero is undefined"
}
```

#### 5. Negative Square Root (Domain Validation Error)
```bash
curl -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" \
  -d '{"a": -9}'
```
**Response (422 Unprocessable Entity):**
```json
{
  "error": "square root of negative number is undefined"
}
```

#### 6. Missing Parameter (Bad Request)
```bash
curl -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" \
  -d '{"a": 5}'
```
**Response (400 Bad Request):**
```json
{
  "error": "field 'b' is required"
}
```

#### 7. Health Check
```bash
curl -i http://localhost:8080/health
```
**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

### Error Handling & Status Codes

| HTTP Status | Condition | Example Response |
|---|---|---|
| `200 OK` | Successful arithmetic evaluation | `{"result": 42}` |
| `400 Bad Request` | Missing operand fields, invalid JSON syntax, empty body, or extraneous payload | `{"error": "field 'a' is required"}` |
| `413 Payload Too Large` | Request body exceeds the 1MB payload limit | `{"error": "request payload exceeds size limit"}` |
| `422 Unprocessable Entity`| Mathematical violation (division by zero, negative square root, undefined value) | `{"error": "division by zero is undefined"}` |
| `500 Internal Server Error`| Unhandled system exception or panic recovered by middleware | `{"error": "internal server error"}` |

---

## Testing & Coverage

The application maintains a comprehensive automated test suite across both stacks with 82 passing tests (50 backend, 32 frontend) and verified coverage metrics.

### Test Execution Summary

| Stack | Runner | Scope | Tests Passed | Status | Statement Coverage |
|---|---|---|---|---|---|
| Backend (Go) | `go test` | 3 packages | 50 / 50 | PASS | >94% (`internal/`) |
| Frontend (React/TS) | Vitest 4.1.11 | 3 test files | 32 / 32 | PASS | Verified suites |

---

### Backend Tests & Coverage Report

```text
=== Package Coverage Breakdown ===
calculator/internal/calculator   95.0% of statements
calculator/internal/handler      91.4% of statements
calculator/internal/server      100.0% of statements
total (internal packages):       >94% of statements
total (including cmd/api):       71.5% of statements

=== Verified Test Suites (50 passed) ===
PASS: TestServiceBinaryOperations (14 subtests: arithmetic, precision, negative zero, limits)
PASS: TestServiceUnaryOperations (5 subtests: sqrt, percentage, zero handling)
PASS: TestHandler (19 subtests: JSON decode, size caps, unprocessable math, bad requests)
PASS: TestRouter (11 subtests: routes, health, CORS headers, preflight OPTIONS)
PASS: TestRecoveryMiddleware (1 subtest: panic recovery and 500 error formatting)
```

#### Run Backend Tests
```bash
cd backend
go test -v -cover ./...
```

Generate coverage profile:
```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

### Frontend Tests Report

```text
=== Verified Test Suites (32 passed) ===
✓ src/utils/format.test.ts (17 tests)
  - Digit buffer limits (16 digits), leading zero replacement
  - Decimal point enforcement, backspace logic, sign toggling
  - IEEE 754 precision smoothing and scientific notation fallback
✓ src/hooks/useCalculator.test.ts (10 tests)
  - State machine initialization, operand buffering, chained binary operations
  - Immediate unary operations, 422 error recovery, history tape persistence
✓ src/components/Calculator.test.tsx (5 tests)
  - Chassis DOM rendering, keypad click handling, keyboard event binding
  - Async API integration, history drawer interaction, and state recall

Test Files  3 passed (3)
     Tests  32 passed (32)
  Duration  1.30s
```

#### Run Frontend Tests
```bash
cd frontend
npm run test:run
```

Run tests in watch mode:
```bash
cd frontend
npm test
```

---

## Design Decisions & Assumptions

1. **Standard Library Over Heavy Frameworks (Go)**:
   - Utilizes Go's standard library `net/http` router with route patterns introduced in Go 1.22 (`POST /api/v1/add`).
   - Keeps binary size minimal, avoids external dependency vulnerabilities, and ensures idiomatic Go architecture.

2. **Layered Clean Architecture**:
   - `internal/calculator`: Domain logic independent of transport protocols.
   - `internal/handler`: HTTP request decoding, schema validation, and response formatting.
   - `internal/server`: Routing, CORS, panic recovery, and logging middleware.
   - `cmd/api`: Composition root handling environment configuration and graceful OS signal shutdown (`SIGINT`, `SIGTERM`).

3. **IEEE 754 Floating-Point Normalization**:
   - Binary floating-point arithmetic can yield inaccuracies (such as `0.1 + 0.2 = 0.30000000000000004`).
   - The backend service normalizes results using `strconv.FormatFloat(val, 'g', 15, 64)`, producing clean decimals and converting `-0` to `0`.

4. **Production Reverse Proxy Architecture**:
   - In containerized deployments, Nginx acts as the single entry point on port 3000, serving static React assets and proxying `/api/` and `/health` requests to the internal Go microservice.
   - Avoids cross-origin overhead in production while local development maintains permissive CORS for `http://localhost:5173`.

5. **Client-Side History with LocalStorage Persistence**:
   - The calculator calculation engine is entirely stateless on the backend.
   - History entries are captured on the client, stored in `localStorage`, and capped at 50 items to ensure fast loads, privacy, and offline persistence.

6. **Accessibility & Usability**:
   - The display uses `aria-live="polite"` so screen readers announce results.
   - High-contrast visual cues distinguish primary operands, operators, and clear actions.
   - The keypad uses CSS Grid to maintain consistent tactile alignment across viewport sizes.
