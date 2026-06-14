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