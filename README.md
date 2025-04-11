# Tarot API

- [Tarot API](#tarot-api)
  - [Overview](#overview)
  - [Features](#features)
  - [Tech Stack](#tech-stack)
  - [Installation](#installation)
    - [Clone the Repository](#clone-the-repository)
    - [Environment Variables](#environment-variables)
  - [Database Setup](#database-setup)
    - [Using Docker (recommended)](#using-docker-recommended)
    - [Manual Setup (if PostgreSQL is installed locally)](#manual-setup-if-postgresql-is-installed-locally)
  - [Usage](#usage)
    - [Build the project](#build-the-project)
    - [Start the server](#start-the-server)
  - [Swagger API Documentation](#swagger-api-documentation)
    - [Generate docs:](#generate-docs)
    - [View in browser:](#view-in-browser)
  - [Demo](#demo)
  - [Testing](#testing)
    - [Run unit tests:](#run-unit-tests)
    - [Run integration tests with Docker:](#run-integration-tests-with-docker)
  - [License](#license)


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

## Demo

A live demo of the Tarot API is available here:

👉 [https://tarot-book.ru/docs/](https://tarot-book.ru/docs/)

This is the automatically generated Swagger UI served from the Go application and proxied through Apache. It reflects the current public API structure.

Please note:
- The documentation includes a simplified version of card interpretations.
- The full database of interpretations (used internally by the Telegram bot) is not exposed in this public API.
- There may be minor cosmetic issues in the URL (e.g. extra slashes in redirects), but all routes work correctly.

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
