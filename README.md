# ⏱️ imitr - High-Performance API Rate Limiter

Limitr is an **API rate limiting service** built in Go, with PostgreSQL for API key management and Redis for low-latency request tracking.
It implements a **Sliding Window Log algorithm**, ensuring fair usage and preventing abuse in multi-tenant systems.

---

## Features

- ⚡ **Sliding Window Rate Limiter** - fast, fair, hot-swappable
- 🔑 **API Key Management** - create, list and delete API keys
- ⚖ **Per-key Plan Solution** - each key has its own assigned plan 
- 📊 **Usage Endpoint** (`/usage`) - real-time usage metrics

---

## Demo

Coming soon!

---

## Tech Stack

- **Go + Chi** - HTTP routing & middleware
- **sqlc** - type-safe, idiomatic code generation from SQL
- **goose** - streamlined database migrations handling
- **PostgreSQL** - API key and client plan storage
- **Redis** - blazingly fast request tracking & caching
- **Docker** - quick, safe and automated deployment

---

## Algorithm choice

Limitr currently implements a **Sliding Window Log**, but the code is structured so a **Token Bucket** or any other rate limiting algorithm could be put in its place.

### Sliding Window (current)

- Requests are tracked precisely within a rolling time window.
- If a client keeps spamming above their limit, they’ll stay locked out until their average request rate drops.

### Token Bucket (possible alternative)

- Requests consume "tokens" from a bucket that refills over time.
- Allows bursts while enforcing an average rate.

---

## Quick start & usage
Prerequisites:
- git
- Docker (engine, compose)
- GNU Make (optional)

```sh
# Clone this repo
git clone https://github.com/33313/limitr.git && cd limitr

# Start the app with GNU Make + Docker
make docker
# Or use `docker compose up --build`

# Create an API key
curl -X POST http://localhost:3000/keys -v
# Example response: '{"api_key": "<YOUR_KEY>"}'

# Use your key to make requests
curl -H "x-api-key: <YOUR_KEY>" http://localhost:3000/usage -v
# Example response (JSON):
# {
#   "limit": 60,
#   "used": 17,
#   "remaining": 43
# }
```

### Adding more routes

1. Edit `./internal/server/routes.go` to add routes and assign handlers
2. Edit `./internal/server/handlers.go` to add handler functions

---

## 📌 Future improvements
- Add alternative rate limit algorithms (token bucket, leaky bucket, fixed window)
- Add Prometheus integration for observability
- Add tests and set up a CI/CD pipeline with GitHub Actions

