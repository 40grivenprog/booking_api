# Booking API

A RESTful API for managing bookings and appointments built with Go, PostgreSQL, and Docker. This project follows a modular architecture pattern with clean separation of concerns.

## Features

- **User Management**: Client and professional user registration and authentication
- **Appointment System**: Book and manage appointments between clients and professionals
- **Type-Safe Database Operations**: Using SQLC for generated, type-safe database queries
- **PostgreSQL Database**: With migrations and proper schema management
- **JWT Authentication**: Secure token-based authentication
- **Structured Logging**: Comprehensive logging with Zerolog
- **Docker Containerization**: Full containerization with Docker & Docker Compose
- **Modular API Architecture**: Clean, maintainable API structure

## Tech Stack

- **Language**: Go 1.23+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: SQLC (type-safe SQL code generation)
- **Migrations**: golang-migrate
- **Authentication**: JWT
- **Logging**: Zerolog
- **Containerization**: Docker & Docker Compose

## Project Structure

```
booking_api/
├── cmd/                           # Application entrypoints
│   └── main.go                   # Main application
├── internal/                      # Private application code
│   ├── api/                      # API layer
│   │   ├── handlers.go           # Main API registration
│   │   ├── clients/              # Client API module
│   │   │   ├── schema.go         # Request/Response schemas
│   │   │   ├── clients_repository.go # Repository interface
│   │   │   ├── handler.go        # Route registration
│   │   │   └── controller.go     # Business logic
│   │   ├── professionals/        # Professional API module
│   │   │   ├── schema.go         # Request/Response schemas
│   │   │   ├── professionals_repository.go # Repository interface
│   │   │   ├── handler.go        # Route registration
│   │   │   └── controller.go     # Business logic
│   │   └── common/               # Shared utilities
│   │       ├── logger.go         # Logger utilities
│   │       └── error_response.go # Error handling
│   ├── config/                   # Configuration management
│   ├── database/                 # Database connection
│   ├── migrations/               # Database migrations
│   └── repository/               # Database queries (SQLC generated)
├── pkg/                          # Public library code
│   └── server/                   # HTTP server setup
├── configs/                      # Configuration files
├── scripts/                      # Database initialization scripts
├── docker-compose.yml            # Docker services
├── Dockerfile                    # Application container
├── Makefile                      # Build and run commands
└── sqlc.yaml                     # SQLC configuration
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd booking_api
```

2. Copy the example configuration:
```bash
cp configs/server-config.yaml.example configs/server-config.yaml
```

3. Start the services:
```bash
make start
```

This will start:
- PostgreSQL database on port 5432
- pgAdmin on port 8081
- The API server on port 8080

### Configuration

Edit `configs/server-config.yaml` to customize:
- Database connection settings
- JWT secret key
- Server port and timeouts
- Logging configuration

## API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### Client Registration
**POST** `/api/clients/register`

Creates a new client user account.

**Request Body:**
```json
{
  "username": "john_doe",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890"
}
```

**Response:**
```json
{
  "message": "Client registered successfully",
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "first_name": "John",
    "last_name": "Doe",
    "user_type": "client",
    "phone_number": "+1234567890",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get Professionals
**GET** `/api/professionals`

Retrieves all professional users.

**Response:**
```json
{
  "professionals": [
    {
      "id": "uuid",
      "username": "dr_smith",
      "first_name": "Dr. Jane",
      "last_name": "Smith",
      "user_type": "professional",
      "phone_number": "+1234567890",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "count": 1
}
```

#### Professional Sign In
**POST** `/api/professionals/sign_in`

Authenticates a professional user.

**Request Body:**
```json
{
  "username": "dr_smith",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Professional signed in successfully",
  "user": {
    "id": "uuid",
    "username": "dr_smith",
    "first_name": "Dr. Jane",
    "last_name": "Smith",
    "user_type": "professional",
    "phone_number": "+1234567890",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Health Check
**GET** `/health`

Check API health status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE,                    -- Telegram chat ID
    username VARCHAR(255) NOT NULL,           -- Username (required)
    first_name VARCHAR(255) NOT NULL,         -- First name (required)
    last_name VARCHAR(255) NOT NULL,          -- Last name (required)
    user_type user_type NOT NULL,             -- 'client' or 'professional'
    phone_number VARCHAR(20),                 -- Optional phone number
    password_hash VARCHAR(255),               -- For professionals only
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Appointments Table
```sql
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type appointment_type NOT NULL,           -- 'appointment' or 'unavailable'
    client_id UUID REFERENCES users(id),     -- Can be null for unavailable
    professional_id UUID NOT NULL REFERENCES users(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status appointment_status DEFAULT 'pending',
    cancellation_reason TEXT,
    cancelled_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Usage

### Database Management

#### Run Migrations
```bash
make migrate-up
```

#### Create New Migration
```bash
make migrate-create
# Enter migration name when prompted
```

#### Generate SQLC Code
```bash
make sqlc-generate
```

### Development Commands

```bash
# Start all services
make start

# Stop all services
make stop

# View logs
make logs

# Build the application
make build

# Restart services
make restart

# Clean up everything
make clean

# Database commands
make db-shell       # Connect to PostgreSQL shell
make pgadmin        # Open pgAdmin in browser

# Migration commands
make migrate-up     # Run migrations up
make migrate-down   # Run migrations down
make migrate-create # Create new migration

# SQLC commands
make sqlc-generate  # Generate SQLC code
make sqlc-validate  # Validate SQLC configuration
```

### Accessing Services

- **API Server**: http://localhost:8080
- **pgAdmin**: http://localhost:8081
  - Email: admin@booking.com
  - Password: admin
- **PostgreSQL**: localhost:5432
  - Database: booking_db
  - User: booking_user
  - Password: booking_pass

## Development

### Adding New Features

1. **Create Database Migrations**: Add new migration files in `internal/migrations/`
2. **Add SQL Queries**: Create query files in `internal/repository/queries/`
3. **Generate SQLC Code**: Run `make sqlc-generate`
4. **Create API Module**: Follow the pattern in `internal/api/`
   - Create schema.go for request/response types
   - Create repository interface
   - Create handler.go for route registration
   - Create controller.go for business logic
5. **Register Routes**: Add to `internal/api/handlers.go`

### API Module Structure

Each API module follows this pattern:

```
module_name/
├── schema.go              # Request/Response schemas
├── module_repository.go   # Repository interface
├── handler.go             # Route registration
└── controller.go          # Business logic
```

### Testing the API

```bash
# Test client registration
curl -X POST http://localhost:8080/api/clients/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890"
  }'

# Test get professionals
curl http://localhost:8080/api/professionals

# Test professional sign in
curl -X POST http://localhost:8080/api/professionals/sign_in \
  -H "Content-Type: application/json" \
  -d '{
    "username": "dr_smith",
    "password": "password123"
  }'

# Test health check
curl http://localhost:8080/health
```

## Architecture

### Modular Design
The API follows a modular architecture where each feature (clients, professionals, etc.) is self-contained with:
- **Schemas**: Request/response type definitions
- **Repository Interface**: Database operation contracts
- **Handler**: Route registration and middleware
- **Controller**: Business logic implementation

### Database Layer
- **SQLC**: Generates type-safe Go code from SQL queries
- **Migrations**: Version-controlled database schema changes
- **PostgreSQL**: Robust, ACID-compliant database

### Error Handling
- Centralized error response handling
- Structured logging with request context
- Proper HTTP status codes

## Security Considerations

- **Password Hashing**: Currently using simple comparison (implement proper hashing for production)
- **JWT Tokens**: Ready for token-based authentication
- **Input Validation**: Request validation using Gin binding
- **SQL Injection**: Prevented by SQLC's type-safe queries

## Contributing

1. Follow the existing code structure and patterns
2. Add tests for new features
3. Update documentation for API changes
4. Use conventional commit messages

## License

[Add your license here]
