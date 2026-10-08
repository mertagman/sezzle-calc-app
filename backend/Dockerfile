FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

FROM alpine:3.21

RUN adduser -D -u 1000 appuser

WORKDIR /app
COPY --from=builder /app/server /app/server

USER appuser

EXPOSE 8080

CMD ["/app/server"]