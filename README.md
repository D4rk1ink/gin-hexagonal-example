# Gin Hexagonal Example

This is a RESTful API built using the [Gin](https://github.com/gin-gonic/gin) web framework in Go. The project follows **Hexagonal Architecture** to ensure a clean separation of concerns and maintainability.

---

## Prerequisites

Before running the project, ensure you have the following installed:

- **Go**: Version 1.24.
- **Docker** and **Docker Compose**: For containerized development.
- **Task**: For task automation.
- **oapi-codegen**: For OpenAPI code generation. (For development)
- **Air**: For hot reloading during development. (For development)

---

## Getting Started

#### Clone the Repository

```bash
git clone https://github.com/D4rk1ink/gin-hexagonal-example
```
#### Run the Application

```bash
cd gin-hexagonal-example
task prod
```
---

## API Documentation

The API is documented using OpenAPI. You can find the specification in the `docs/server/doc.yaml` file or enter to [http://localhost:8080/swagger](http://localhost:8080/swagger) or [Postman Collection](https://github.com/D4rk1ink/gin-hexagonal-example/blob/master/postman_collection.json). To view the documentation:

#### Example API

##### Register API
```bash
curl --location 'http://localhost:8080/api/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "FirstName LastName",
    "email": "user@email.com",
    "password": "password",
    "confirm_password": "password"
}'
```
##### Login API
```bash
curl --location 'http://localhost:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user@email.com",
    "password": "password"
}'
```
##### Create User API
Note: Please adjust Authorization header
```bash
curl --location 'http://localhost:8080/api/users' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data-raw '{
    "name": "FirstName2 LastName2",
    "email": "user2@email.com",
    "password": "password",
    "confirm_password": "password"
}'
```
##### Get Users API
Note: Please adjust Authorization header
```bash
curl --location 'http://localhost:8080/api/users' \
--header 'Authorization: ••••••'
```
##### Get User By Id API
Note: Please adjust Authorization header and user id
```bash
curl --location 'http://localhost:8080/api/users/681baad75bf4df5c5f1b5d9b' \
--header 'Authorization: ••••••'
```
##### Update User API
Note: Please adjust Authorization header and user id
```bash
curl --location --request PATCH 'http://localhost:8080/api/users/681baf295bf4df5c5f1b5d9c' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data-raw '{
    "name": "NewFirstName LastName",
    "email": "new_user@email.com"
}'
```
##### Delete User API
Note: Please adjust Authorization header and user id
```bash
curl --location --request DELETE 'http://localhost:8080/api/users/681baf295bf4df5c5f1b5d9c' \
--header 'Authorization: ••••••'
```
---
## API Logging
Below is an example of a log entry generated for an API request:

```json
{
  "level": "info",
  "ts": "2025-05-07T20:08:34Z",
  "caller": "logger/logger.go:39",
  "msg": "",
  "correlation_id": "cb3c0d95-7251-45e9-8c4d-0d123efc24a9",
  "method": "POST",
  "path": "/api/auth/register",
  "status": 409,
  "duration": 0.001125514
}
```

#### Log Fields Explained:
- **`level`**: The severity level of the log (e.g., `info`, `error`).
- **`ts`**: The timestamp of the log entry in RFC3339.
- **`caller`**: The file and line number where the log was generated.
- **`msg`**: A short message describing the log event.
- **`correlation_id`**: A unique identifier for tracing the request across services.
- **`method`**: The HTTP method of the request (e.g., `GET`, `POST`).
- **`path`**: The API endpoint that was accessed.
- **`status`**: The HTTP status code returned by the server.
- **`duration`**: The time taken to process the request (in seconds).

---

## Project Structure

```plaintext
├── cmd/
│   └── app/
│       └── main.go         # Main application bootstrap where the server is initialized and started
├── config/                 # Configuration folder containing environment variables and settings
├── container/
│   ├── Dockerfile          # Dockerfile for building the image
│   └── docker-compose.yml
├── docs/
│   └── server/
│       └── doc.yaml        # OpenAPI specification for the API endpoints
├── integration/            # Integration tests to validate the interaction between components
├── internal/
│   ├── application/        # Application layer containing client interface
│   │   └── handler/
│   │       ├── http/       # HTTP-specific handlers for routing and request processing
│   │       └── scheduler/  # Background job schedulers for periodic tasks
│   ├── core/
│   │   ├── domain/         # Domain models and entities representing the business logic
│   │   ├── dto/            # Data Transfer Objects for request and response payloads
│   │   ├── error/          # Custom error definitions
│   │   ├── port/           # Interfaces for dependency inversion (e.g., repositories, services)
│   │   └── service/        # Business logic implementations for the core domain
│   ├── infrastructure/     # Infrastructure-specific implementations and utilities
│   │   ├── config/         # Configuration management and environment setup
│   │   ├── database/
│   │   ├── dependency/     # Dependency injection setup for the application
│   │   ├── hash/
│   │   ├── jwt/
│   │   ├── logger/
│   │   └── repository/     # Datasource repositories for external data
│   └── util/               # Utility functions and helpers used across the application
├── go.mod
├── go.sum
├── README.md               # Project documentation
└── Taskfile.yml            # Task runner configuration for automating common tasks
```
---

## Development
### Install Dependencies

```bash
go mod tidy
```

### Run the Application (Development Mode)

```bash
task dev
```

### Running Unit Tests

Run the following command to execute the test suite:

```bash
task utest
```
Watch mode
```bash
task utest-watch
```
With coverage
```bash
task utest-coverage
```
### Running Integration Tests

Run the following command to execute the integration test:

```bash
task dev
task itest
```
Watch mode
```bash
task dev
task itest-watch
```

## Acknowledgments

- [Hexagonal Architecture](https://www.baeldung.com/hexagonal-architecture-ddd-spring)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Air](https://github.com/cosmtrek/air)
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [Docker](https://www.docker.com/)
- [Task](https://github.com/go-task/task)