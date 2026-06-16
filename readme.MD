# Go User API — Backend Development Task

A RESTful API built in Go that manages users with a name and date of birth, calculating each user's current age dynamically on every fetch.

## Tech Stack

- **Go** 1.21
- **GoFiber** — HTTP web framework
- **PostgreSQL + SQLC** — type-safe, generated database access layer
- **Uber Zap** — structured JSON logging
- **go-playground/validator** — request input validation

## Project Structure

```
.
├── cmd/server/main.go        # Application entry point
├── config/                   # Environment configuration loading
├── db/
│   ├── migrations/           # Raw SQL migration files
│   ├── schema/                # Schema used by SQLC for type generation
│   └── queries/               # SQL queries used by SQLC
├── internal/
│   ├── handler/               # HTTP request/response handling
│   ├── service/                # Business logic (age calculation, orchestration)
│   ├── db/                     # SQLC-generated database access code
│   ├── database/               # Database connection setup
│   ├── routes/                  # Route registration
│   ├── middleware/              # Request ID + request duration logging
│   ├── models/                  # Shared data structures
│   └── logger/                  # Zap logger wrapper
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Setup & Run Instructions

### Prerequisites
- Go 1.21+
- Docker Desktop (for running PostgreSQL) — or an existing local PostgreSQL instance

### 1. Clone the repository and install dependencies
```bash
git clone <your-repo-url>
cd go-user-api
go mod tidy
```

### 2. Configure environment variables
```bash
cp .env.example .env
```
Update `.env` if needed. Defaults used in development:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_api_db
DB_DRIVER=postgres
SERVER_PORT=3000
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
```

### 3. Start PostgreSQL
```bash
docker run --name user-api-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=user_api_db \
  -p 5432:5432 -d postgres:15
```

### 4. Apply the database migration
```bash
docker exec -i user-api-db psql -U postgres -d user_api_db < db/migrations/001_create_users_table.sql
```

### 5. (Optional) Regenerate SQLC code
Only needed if you modify the SQL queries in `db/queries/`:
```bash
sqlc generate
```

### 6. Run the server
```bash
go run cmd/server/main.go
```
The server starts on `http://localhost:3000`.

## API Endpoints

| Method | Endpoint     | Description                          |
|--------|--------------|---------------------------------------|
| POST   | `/users`     | Create a new user                     |
| GET    | `/users/:id` | Get a user by ID, with calculated age |
| PUT    | `/users/:id` | Update a user's name/DOB              |
| DELETE | `/users/:id` | Delete a user (returns 204)           |
| GET    | `/users`     | List all users, with calculated ages  |

### Example: Create a user
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

### Example: Get a user (age calculated dynamically)
```bash
curl http://localhost:3000/users/1
```
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 36
}
```

## Validation & Error Handling
Input is validated using `go-playground/validator` (required name, valid date format for `dob`). Invalid requests return a `400 Bad Request` with a descriptive error message instead of a server crash.

## Logging
All key actions (server startup, database connection, user creation/update/deletion, errors) are logged in structured JSON format using Uber Zap.

## Middleware
- **Request ID**: every response includes an `X-Request-Id` header for traceability.
- **Request duration logging**: each request's processing time is logged.

## Testing
Unit tests cover the age calculation logic, including edge cases (leap years, birthdays not yet reached this year):
```bash
go test -v ./...
```

## Docker
A `Dockerfile` and `docker-compose.yml` are included for containerized deployment of the application alongside PostgreSQL.

## Notes on Design Choices
- **SQLC** was used instead of an ORM or raw inline queries to get compile-time type safety on database access while keeping full control over the SQL itself.
- **Age** is never stored in the database — it's calculated on every read using Go's `time` package, since DOB is the only source of truth and age changes daily.
