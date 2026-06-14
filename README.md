# Cinema Booking

Simple seat booking app written in Go with a Redis-backed hold flow and a static browser UI.

## Features

- Browse movies and available seats
- Hold a seat for a short time
- Confirm or release a held seat
- Redis-backed session storage for seat holds

## Project Structure

- `cmd/main.go` - HTTP server entrypoint
- `internal/booking` - booking domain, service, handlers, and stores
- `internal/booking/adapters/redis` - Redis client helper
- `static/index.html` - browser UI
- `docker-compose.yaml` - local Redis and Redis Commander setup

## Requirements

- Go 1.22+
- Redis if you want to run the Redis-backed flow
- Docker Desktop if you want the provided local Redis stack

## Environment

Copy [.env.example](.env.example) to `.env` if you want to keep local settings in one place.

### Using a free hosted Redis (Upstash)

You can use Upstash's free Redis plan for development. Steps:

1. Create a free account at https://upstash.com and create a Redis database.
2. Copy the REST/Redis endpoint (host:port) and the password.
3. In Render (or your host) set environment variables:

```text
REDIS_ADDR=<host:port>
REDIS_PASSWORD=<your-password>
```

4. Paste only the Redis host or URI into the variable value. Do not paste the full `redis-cli --tls -u ...` command.
5. The app reads `REDIS_ADDR` and `REDIS_PASSWORD` automatically when starting, and also accepts `REDIS_URL` if you prefer a single connection string.

Upstash offers a small free tier suitable for development and light testing.

## Run Locally

### With Docker

Start Redis and Redis Commander:

```powershell
docker compose up -d
```

Run the app:

```powershell
go run ./cmd/main.go
```

Open the app in your browser:

```text
http://localhost:8080
```

### Without Docker

You can also run the app without the compose stack if you switch to the in-memory store in the code. The current default flow expects Redis.

## Tests

Run all tests:

```powershell
go test ./...
```

## Notes

- Seat holds expire automatically after a short TTL.
- Closing the tab triggers a best-effort release, but TTL is still the fallback.