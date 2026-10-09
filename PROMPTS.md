# AI Collaboration Log (PROMPTS.md)

## Summary of Usage
- **Tools**: Gemini / Claude
- **Scope**: Architecture scaffolding, test troubleshooting, CSS refinements, and Docker Compose orchestration.
- **Workflow**: Iterative design review, manual code auditing, and independent test verification.

---

## 1. Project Scaffolding & Full-Stack Architecture
- **Objective**: Establish a clean separation of concerns with a modular Go REST API and a Vite React SPA.
- **Key Prompt(s)**:
  > "Design an idiomatic project layout for a full-stack calculator application: Go backend using standard library routing (`net/http`) under `internal/` (service, handler, server) and React/TypeScript frontend under `src/`."
- **Outcome & Human Review**: Adopted a clean architecture separating calculation domain logic (`internal/calculator`) from HTTP request validation (`internal/handler`) and routing/middleware (`internal/server`).

---

## 2. Backend Service, Handler & Unit Testing
- **Objective**: Implement arithmetic operations, input validation, and high test coverage using Go table-driven unit tests.
- **Key Prompt(s)**:
  > "Implement the calculator domain service supporting binary and unary operations (+, -, *, /, pow, sqrt, percent) with strict error handling for division by zero and negative square roots. Write table-driven unit tests covering edge cases, precision limits, and achieve >90% test coverage."
  > "Create HTTP handlers for /api/v1/* routes with JSON request decoding, standardized error responses (400 Bad Request, 422 Unprocessable Entity), and handler unit tests using `net/http/httptest`."
- **Outcome & Human Review**: Validated business logic via `service_test.go`, verified HTTP status codes and panic recovery middleware via `handler_test.go` and `router_test.go`, achieving verified coverage across all backend packages.

---

## 3. Frontend Architecture & Testing
- **Objective**: Manage calculator state, history tape persistence, and resolve JSDOM/Vitest test suite regressions.
- **Key Prompt(s)**:
  > "Create a custom `useCalculator` hook managing input buffer, operations, and localStorage history. Write comprehensive Vitest integration tests covering keypad clicks, keyboard inputs, and mock isolation."
  > "Fix Vitest test failures: resolve 'localStorage.clear is not a function' under Node 22, and eliminate TestingLibrary query collisions for identical numbers in the display and history panel."
- **Outcome & Human Review**: Implemented a standalone Storage mock in `setupTests.ts`, replaced ambiguous text queries with explicit `data-testid` and `getByRole` selectors, and enforced `mockResolvedValueOnce` for clean test isolation.

---

## 4. UI, Layout & Accessibility
- **Objective**: Build a responsive retro-hardware calculator interface with CSS Grid, dual-line display, history tape, and full ARIA support.
- **Key Prompt(s)**:
  > "Design a retro-hardware desktop calculator in React using CSS Grid for keypad layout, a dual-line LCD display with tabular-nums, a persistent history tape drawer, and complete ARIA/keyboard accessibility."
- **Outcome & Human Review**: Implemented modular components (Calculator, Display, Keypad, History), resolved button shifting via rigid grid constraints, and mapped physical keyboard events to accessible DOM nodes.

---

## 5. Docker Orchestration & Production Reverse Proxy
- **Objective**: Containerize both services and unite them under a single root orchestration setup with automated health checks.
- **Key Prompt(s)**:
  > "Write multi-stage Dockerfiles for Go (minimal Alpine runner) and React (Vite build + Nginx 1.27 Alpine). Configure Nginx reverse proxy using template substitution (`envsubst`) to route `/api/` and `/health` to the backend."
  > "Create a root `docker-compose.yml` exposing the frontend on port 3000, enforcing dependency order via backend healthchecks (`service_healthy`), and isolating communication inside a shared bridge network."
- **Outcome & Human Review**: Replaced non-standard port 5173 with standard port 3000, validated Nginx dynamic upstream proxying to `http://backend:8080`, and verified end-to-end execution via curl checks.