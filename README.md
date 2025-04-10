# Tarot API

## Overview
Tarot API is a RESTful service that provides access to a structured database of classical Tarot data, including decks, spreads, interpretations (meanings), and card metadata.

This project is suitable for building Tarot bots, educational tools, divination assistants, and web/mobile Tarot clients.

## Features

- Full CRUD for:
  - Decks
  - Cards (Major & Minor Arcana)
  - Sources
  - Spreads
  - Suits & Ranks
  - Meanings (Major & Minor Arcana)
- Meaning filters (e.g. by `source`, `position`, `number`, `suit`)
- Swagger UI documentation
- JSON API responses
- Integration test suite using isolated PostgreSQL
- Modular, idiomatic Go codebase

## Tech Stack

- **Go** 1.21+
- **PostgreSQL**
- **Echo** web framework
- **Swagger / Swaggo** for API docs
- **Docker Compose** for local dev/test DB
- **Testify** for testing

---

## Installation

### Clone the Repository

```sh
git clone https://github.com/ilbagatto/tarot-api.git
cd tarot-api
```

### Environment Variables

Create a `.env` file:

```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5442
DB_USER=tarot
DB_PASSWORD=yourpassword
DB_NAME=tarot
POSTGRES_DSN=postgresql://tarot:yourpassword@localhost:5442/tarot?sslmode=disable

# Logging format: use "color", "development", or "json"
LOG_FORMAT=development

# Public URL where card images are served from
BASE_URL=https://yourdomain.com/static
```

---

## Database Setup

### Using Docker (recommended)

```sh
docker-compose up -d
```

For testing environment:
```sh
docker-compose -f docker-compose.test.yml up -d
```

### Manual Setup (if PostgreSQL is installed locally)

```sh
createdb tarot
psql -U tarot -d tarot -f setup-db/init.sql
```

---

## Usage

### Build the project

```sh
make build
```

### Start the server

```sh
make run
```

API will be available at: [http://localhost:8080](http://localhost:8080)

---

## Swagger API Documentation

### Generate docs:

```sh
make docs
```

### View in browser:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Testing

### Run unit tests:

```sh
make test-unit
```

### Run integration tests with Docker:

```sh
make test-integration
```

---

## License

This project is licensed under the MIT License.

⚠️ Interpretations are demo-only. Real data is omitted for licensing purposes.
