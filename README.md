# Go Gin Project

This is a RESTful API built using the [Gin](https://github.com/gin-gonic/gin) web framework in Go. The project follows **Hexagonal Architecture** to ensure a clean separation of concerns and maintainability.

---

## Features

- **Hexagonal Architecture**: Clear separation of core business logic, adapters, and infrastructure.
- **Gin Framework**: High-performance HTTP web framework.
- **OpenAPI Integration**: API documentation and code generation using `oapi-codegen`.
- **Hot Reloading**: Enabled with [Air](https://github.com/cosmtrek/air) for development.
- **MongoDB**: Database integration with Dockerized MongoDB.
- **Taskfile**: Simplified task automation for common operations.

---

## Prerequisites

Before running the project, ensure you have the following installed:

- **Go**: Version 1.20 or later.
- **Docker** and **Docker Compose**: For containerized development.
- **oapi-codegen**: For OpenAPI code generation.
- **Air**: For hot reloading during development (optional).
- **Task**: For task automation.

---

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/D4rk1ink/gin-hexagonal-example
cd gin-hexagonal-example
```

### Install Dependencies

```bash
go mod tidy
```

### Run the Application (Development Mode)

```bash
task dev
```


---

## Project Structure

```plaintext
├── cmd/                # Application entry points
│   └── main.go         # Main application bootstrap
├── internal/           # Core application logic
│   ├── adapter/        # Adapters for external systems (e.g., HTTP, DB)
│   │   ├── handler/    # HTTP handlers
│   │   └── repository/ # Database repositories
│   ├── application/    # Application services (business logic)
│   ├── domain/         # Core domain models and interfaces
│   └── config/         # Configuration management
├── docs/               # OpenAPI specification
│   └── server/         
│       └── doc.yaml    # API documentation files
├── container/          # Docker Compose files
│   ├── Dockerfile.dev  # Dockerfile for development
│   └── docker-compose.dev.yaml
├── .air.toml           # Air configuration for hot reloading
├── Taskfile.yml        # Task runner configuration
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── README.md           # Project documentation
```

---

## API Documentation

The API is documented using OpenAPI. You can find the specification in the `docs/server/doc.yaml` file or enter to [http://localhost:8080/swagger](http://localhost:8080/swagger). To view the documentation:



---

## Development

### Running Tests

Run the following command to execute the test suite:

```bash
task test
```

### Linting

Ensure your code adheres to Go standards by running:

```bash
task lint
```

---

## Acknowledgments

- [Hexagonal Architecture](https://www.baeldung.com/hexagonal-architecture-ddd-spring)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Air](https://github.com/cosmtrek/air)
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [Docker](https://www.docker.com/)
- [Task](https://github.com/go-task/task)