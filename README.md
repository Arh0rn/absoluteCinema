# absoluteCinema

absoluteCinema is a CRUD web application for managing movie data, developed using Golang with Clean Architecture principles. It utilizes PostgreSQL for persistent storage, JWT tokens for authorization, and exposes a RESTful API documented via Swagger.

## Features
- **Clean Architecture**: Ensures separation of concerns, making the application maintainable and testable.
- **JWT Authentication**: Secure user authentication using JSON Web Tokens.
- **PostgreSQL**: Reliable and scalable relational database storage.
- **REST API**: Easy-to-use API endpoints developed using the Mux router.
- **Swagger Documentation**: Automatically generated API documentation for clarity and ease of integration.

## Installation
go version: **1.23**

Clone the repository:

```bash
git clone <repository_url>
cd absoluteCinema
```

Install dependencies:

```bash
go mod tidy
```

## Environment Setup

Create a `.env` file in the project root with the following structure:

```env
POSTGRES_DB_HOST=localhost
POSTGRES_DB_PORT=5432
POSTGRES_DRIVER_NAME=postgres
POSTGRES_DB_USER=postgres
POSTGRES_DB_PASSWORD=postgres
POSTGRES_DB_NAME=absoluteCinema
POSTGRES_DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

HASH_SALT=salt
JWT_SECRET=secret
```

## Usage

### Build Application

```bash
make build
```

### Run Application

```bash
make run
```

### Swagger Documentation

Generate Swagger documentation:

```bash
make swag
```

Access Swagger UI:

```
http://localhost:8080/swagger/index.html
```

### Database Migrations

Install migration tool:

```bash
make install-migrate
```

Apply migrations:

```bash
make migrate-up
```

Revert migrations:

```bash
make migrate-down
```

## Cleaning and Rebuilding

Clean the build:

```bash
make clean
```

Rebuild the application:

```bash
make rebuild
```

## Printing Environment Variables

To verify loaded environment variables:

```bash
make print.env
```

## Swagger JSON Endpoint

Swagger documentation JSON is available at:

```bash
http://<HOST>:<PORT>/swagger/doc.json
```

Replace `<HOST>` and `<PORT>` with your configured values.

