# Apartment Management Backend

Go backend service for the Apartment Management system.

## Status

Backend setup instructions are ready. API route/path documentation is intentionally not included yet because that part is still in progress.

## Prerequisites

- Go 1.25+
- PostgreSQL
- Docker and Docker Compose (optional, for container-based setup)

## Environment Setup

1. Copy `.env.example` to `.env`.
2. Configure the required values:
   - `PORT`
   - `JWT_SECRET`
   - `DB_HOST`
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`
   - `DB_PORT`
   - `POSTGRES_USER`
   - `POSTGRES_PASSWORD`
   - `POSTGRES_DB`

## Install Dependencies

```bash
go mod tidy
```

## Run Locally

```bash
go run main.go
```

## Run with Docker

```bash
docker compose up --build
```

## Run Tests

```bash
go test ./...
```
xd