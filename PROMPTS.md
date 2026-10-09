




calculator-app/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Server bootstrap and route registration
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── service.go           # Arithmetic logic and domain errors
│   │   │   └── service_test.go      # Table-driven unit tests
│   │   └── handler/
│   │       ├── handler.go           # HTTP request parsing, JSON formatting
│   │       └── handler_test.go      # HTTP response/error tests
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Display.tsx          # Current input, calculation expression, and results
│   │   │   └── Keypad.tsx           # Button layout and grid
│   │   ├── services/
│   │   │   └── api.ts               # HTTP client calling the backend
│   │   ├── types/
│   │   │   └── index.ts             # TypeScript interfaces for request/response payloads
│   │   ├── App.tsx                  # Root calculator state and orchestration
│   │   ├── App.test.tsx             # Frontend interaction and render tests
│   │   ├── index.css                # Base styling
│   │   └── main.tsx                 # DOM entry point
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── Dockerfile
├── docker-compose.yml               # Local multi-container runner
├── PROMPTS.md                       # AI prompt log required by Sezzle
└── README.md                        # Architecture, setup guide, and API documentation



frontend/
├── src/
│   ├── api/
│   │   ├── client.ts            # Backend API HTTP client
│   │   └── types.ts             # API request and response interfaces
│   ├── components/
│   │   ├── Calculator.test.tsx  # End-to-end component integration test
│   │   ├── Calculator.tsx       # Main container connecting hook to UI
│   │   ├── Display.tsx          # Numerical display and error banner
│   │   └── Keypad.tsx           # Button grid and click handlers
│   ├── hooks/
│   │   ├── useCalculator.test.ts # State machine and business logic unit tests
│   │   └── useCalculator.ts      # State management and API execution logic
│   ├── App.tsx                  # Root application component
│   ├── index.css                # Global and component layout styling
│   ├── main.tsx                 # React DOM mount point
│   └── setupTests.ts            # Testing library matchers configuration
├── Dockerfile                   # Multi-stage production container build
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.node.json
└── vite.config.ts

go run ./cmd/api

go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

docker build -t calculator-backend .
docker run -d -p 8080:8080 -e ALLOWED_ORIGIN="http://localhost:5173" --name calc-backend calculator-backend
docker stop calc-backend && docker rm calc-backend




slog

naming conventions


checking the best practices