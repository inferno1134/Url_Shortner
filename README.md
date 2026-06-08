# URL Shortener Notes

## Project Overview

This repository is a simple URL shortening service written in Go. It supports creating a short code for a long URL and redirecting requests from the short code back to the original URL.

## Layers and Responsibilities

1. Service Layer
   - `service/shortner.go`
   - Defines the `Store` interface with `Save(longUrl string) string` and `Get(shortCode string) (string, bool)`.
   - Contains `ShortnerService`, which is the business logic layer and delegates storage operations to the configured store.

2. Store Layer
   - `store/memory.go`
   - `store/postgres.go`
   - `MemoryStore` stores URLs in memory using a map and generates short codes with base62 encoding from an auto-increment counter.
   - `PostgresStore` stores URLs in PostgreSQL and also uses base62 encoding for generated IDs.

3. Handler Layer
   - `handlers/shorten.go`
   - `handlers/redirect.go`
   - `ShortenUrl` handles POST `/shorten`, decoding JSON with a `url` field and returning a shortened URL.
   - `RedirectURL` handles GET `/{code}` and redirects to the stored long URL.

## Main Application Flow

- `main.go` sets up the service and handlers.
- It checks for a `DATABASE_URL` environment variable:
  - If present, it uses `PostgresStore`.
  - Otherwise, it uses `MemoryStore`.
- It exposes two routes:
  - `POST /shorten`
  - `GET /{shortCode}`

## Encoding

- Short codes are generated using base62 encoding.
- `MemoryStore` converts an integer counter into a base62 string.
- `PostgresStore` assigns a numeric ID from the database and converts it into base62.

## API

### Create shortened URL

POST `/shorten`

Request body:

```json
{
  "url": "https://example.com"
}
```

Response body:

```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

### Redirect

GET `/{shortCode}`

- Redirects to the original URL if the code exists.

## Notes and Current Behavior

- The current implementation stores data in memory by default.
- If `DATABASE_URL` is configured, it uses PostgreSQL.
- There is no duplicate URL detection; each request generates a new short code.
- The returned short URL uses `http://localhost:8080`.
- No authentication or rate limiting is implemented.

## Run Locally

```sh
go run main.go
```

## Run in Docker

```sh
docker build -t url-shortner .
docker run -p 8080:8080 url-shortner
```
