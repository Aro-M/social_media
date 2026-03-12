# Social Media API

A high-performance backend engine built with Go and PostgreSQL.

### Repository Interface Functions

- **CreateUser**: Inserts a new user record into the database.
- **GetUserByID**: Retrieves a user entity by its unique UUID.
- **GetUserByEmail**: Finds a user by their email address.
- **SearchUsers**: Filters users by first name, last name, and age, excluding the current user.
- **UpdateUser**: Persists modifications to an existing user record.
- **CreateFriendRequest**: Initializes a new friendship or request record.
- **GetFriendRequest**: Retrieves a specific friend request by its ID.
- **UpdateFriendRequest**: Updates the status (pending, accepted, declined) of a request.
- **ListFriendRequests**: Lists all pending friend requests received by a user.
- **ListFriends**: Retrieves all users who have an accepted friendship with the subject.
- **RemoveFriend**: Deletes a friendship or request record.
- **GetFriendshipByParticipants**: Finds a friendship record between two specific user IDs.
- **Close**: Safely terminates the database connection pool.

## Architecture Overview

This project follows Clean Architecture principles:

1.  Controller Layer: Handles HTTP requests and responses (Echo Framework).
2.  Service Layer: Contains business logic.
3.  Repository Layer: Interacts with the database (PostgreSQL).
4.  Entity Layer: Core data models.

## Quick Start

### Option 1: Using Docker (Recommended)

1.  Start Services:
    ```bash
    make up
    ```
2.  Run Migrations:
    ```bash
    make migrate-up
    ```
3.  Access API:
    - Server: http://localhost:8089
    - Documentation: http://localhost:8089/swagger/index.html

### Option 2: Running Locally (Without Docker)

1.  Database Setup:
    - Ensure PostgreSQL is running. You can start just the database in Docker:
      ```bash
      docker-compose up -d db
      ```
2.  Configure Environment:
    - Open .env and set DB_HOST=localhost and DB_PORT=5454 (if using Docker database).
3.  Run Migrations:
    - If using Docker DB, you can still run:
      ```bash
      make migrate-up
      ```
4.  Install Dependencies:
    ```bash
    go mod tidy
    ```
5.  Run the application:
    ```bash
    go run cmd/main.go
    ```

## Makefile Commands

- make up: Start containers.
- make down: Stop containers.
- make test: Run all tests.
- make test-coverage: Run tests and show coverage report.
- make swag: Update Swagger docs.
- make migrate-up: Apply migrations (Docker).

## Features

- User: Register, Login, Update Profile, Search.
- Friends: Send/Accept/Decline Requests, List Friends, Remove Friend.
- Auth: JWT-based authentication.
- Documentation: Auto-generated Swagger UI.

## Testing

```bash
make test
```
- Unit Tests: Test logic in Service and Repository.

### Test Coverage

To see the test coverage report in the terminal:
```bash
make test-coverage
```

To see the detailed coverage report in your browser:
```bash
make test-coverage-html
```
