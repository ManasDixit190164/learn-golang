# Go URL Shortener

Production-style URL shortener backend using Go, PostgreSQL, JWT auth, refresh tokens, clean architecture, migrations, and Docker Compose.

## Features

- Signup/Login
- JWT access token authentication
- Refresh token rotation with hashed refresh tokens in PostgreSQL
- Create short URLs
- Optional custom alias
- List/get/update/deactivate user's URLs
- Public redirect endpoint
- Click tracking
- Basic analytics
- PostgreSQL migrations
- Docker Compose
- Clean layered structure

## Run Locally

```bash
cp .env.example .env
docker compose up -d postgres
make migrate-up
make run
```

Server runs at:

```txt
http://localhost:8080
```

## API Test Flow

### Signup

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Manas","email":"manas@example.com","password":"Password@123"}'
```

Copy `access_token` from the response.

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manas@example.com","password":"Password@123"}'
```

### Create Short URL

```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"original_url":"https://example.com/some/very/long/url","title":"Example"}'
```

### Create Custom Alias

```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"original_url":"https://example.com","custom_alias":"my-link","title":"My Link"}'
```

### Redirect

```bash
curl -i http://localhost:8080/my-link
```

### List URLs

```bash
curl http://localhost:8080/api/v1/urls \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Analytics

```bash
curl http://localhost:8080/api/v1/urls/URL_ID/analytics \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Folder Structure

```txt
cmd/api/main.go               App entry point
internal/config               Environment config
internal/database             PostgreSQL connection
internal/domain               Core entities and DTOs
internal/repository           Database layer
internal/service              Business logic
internal/handler              HTTP handlers
internal/middleware           Auth middleware
internal/router               Routes
internal/utils                JWT, hashing, slug, validation helpers
pkg/response                  Standard API response helper
migrations                    SQL migrations
```

## Notes

For production, add HTTPS, stricter CORS, request rate limiting, refresh-token device sessions, structured logging, background analytics queue, and deployment-specific secrets management.
