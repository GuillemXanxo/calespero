# Calespero User Authentication Service

A Go-based user authentication service with PostgreSQL database integration, designed to be deployed on Render.

## Features

- User registration and authentication
- JWT-based session management
- PostgreSQL database integration
- Hexagonal architecture
- Concurrent-safe database operations

## Prerequisites

- Go 1.19 or later
- PostgreSQL database (hosted on Render)

## Environment Variables

The following environment variables need to be set:

```
DATABASE_URL=postgres://username:password@host:port/database_name
JWT_SECRET_KEY=your_jwt_secret_key
```

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   ├── ports/
│   │   └── services/
│   ├── handlers/
│   └── repositories/
│       └── postgres/
├── pkg/
│   └── auth/
├── templates/
│   ├── login.html
│   ├── new_user.html
│   └── start.html
└── migrations/
    └── 001_create_users_table.sql
```

## Setup

1. Clone the repository
2. Set up the required environment variables
3. Run the database migrations
4. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

```bash
go run cmd/api/main.go
```

The server will start on port 3500.

## API Endpoints

- `GET /login` - Serves the login form
- `POST /login` - Authenticates user and creates session
- `GET /new_user` - Serves the registration form
- `POST /new_user` - Creates a new user
- `GET /start` - Protected route, shows logged-in user's page

## Database Schema

### Users Table
- `id` (VARCHAR) - Primary key
- `email` (VARCHAR) - Unique, required
- `password` (VARCHAR) - Hashed password
- `phone` (VARCHAR) - Required
- `created_at` (TIMESTAMP) - Creation timestamp
- `last_connection` (TIMESTAMP) - Last login timestamp

## Security Features

- Password hashing using bcrypt
- JWT-based authentication
- HTTP-only cookies for token storage
- Concurrent-safe database operations using mutexes 