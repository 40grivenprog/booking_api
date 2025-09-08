# Booking API

A comprehensive RESTful API for managing appointments and bookings between clients and professionals. Built with Go, PostgreSQL, and Docker following clean architecture principles.

## Features

- **User Management**: Client and professional registration and authentication
- **Appointment System**: Complete booking flow with status management
- **Availability Management**: Professional availability checking and time slot management
- **Cancellation System**: Both client and professional appointment cancellation
- **Unavailable Periods**: Professional can mark themselves as unavailable
- **Type-Safe Database Operations**: Using SQLC for generated, type-safe database queries
- **PostgreSQL Database**: With migrations and proper schema management
- **Docker Containerization**: Full containerization with Docker & Docker Compose
- **Modular API Architecture**: Clean, maintainable API structure

## Tech Stack

- **Language**: Go 1.23+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: SQLC (type-safe SQL code generation)
- **Migrations**: golang-migrate
- **Containerization**: Docker & Docker Compose

## Quick Start

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

## API Documentation

### Base URL
```
http://localhost:8080
```

## Client Endpoints

### 1. Client Registration
**POST** `/api/clients/register`

Register a new client account.

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "chat_id": 123456789,
  "phone_number": "+1234567890"
}
```

**cURL:**
```bash
curl -X POST "http://localhost:8080/api/clients/register" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "chat_id": 123456789,
    "phone_number": "+1234567890"
  }'
```

### 2. Get Client Appointments
**GET** `/api/clients/{id}/appointments`

Get all future appointments for a client with optional status filtering.

**Query Parameters:**
- `status` (optional): `pending`, `confirmed`, `cancelled`, `completed`

**cURL:**
```bash
# Get all future appointments
curl -X GET "http://localhost:8080/api/clients/28c31a08-f740-440e-a161-6c8136478e2b/appointments" \
  -H "Content-Type: application/json"

# Get only confirmed appointments
curl -X GET "http://localhost:8080/api/clients/28c31a08-f740-440e-a161-6c8136478e2b/appointments?status=confirmed" \
  -H "Content-Type: application/json"
```

### 3. Cancel Appointment (Client)
**PATCH** `/api/clients/{id}/appointments/{appointment_id}/cancel`

Cancel an appointment as a client.

**Request Body:**
```json
{
  "cancellation_reason": "Need to reschedule"
}
```

**cURL:**
```bash
curl -X PATCH "http://localhost:8080/api/clients/28c31a08-f740-440e-a161-6c8136478e2b/appointments/71a738d8-6695-4fa3-b68a-c58797801258/cancel" \
  -H "Content-Type: application/json" \
  -d '{
    "cancellation_reason": "Need to reschedule"
  }'
```

## Professional Endpoints

### 1. Get All Professionals
**GET** `/api/professionals`

Get list of all professionals.

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/professionals" \
  -H "Content-Type: application/json"
```

### 2. Professional Sign In
**POST** `/api/professionals/sign_in`

Authenticate a professional user.

**Request Body:**
```json
{
  "username": "dr_smith",
  "password": "password123",
  "chat_id": 123456789
}
```

**cURL:**
```bash
curl -X POST "http://localhost:8080/api/professionals/sign_in" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "dr_smith",
    "password": "password123",
    "chat_id": 123456789
  }'
```

### 3. Get Professional Appointments
**GET** `/api/professionals/{id}/appointments`

Get all future appointments for a professional with optional status filtering.

**Query Parameters:**
- `status` (optional): `pending`, `confirmed`, `cancelled`, `completed`

**cURL:**
```bash
# Get all future appointments
curl -X GET "http://localhost:8080/api/professionals/7c065dd1-22b9-4bed-82e2-be973cb6ea47/appointments" \
  -H "Content-Type: application/json"

# Get only pending appointments
curl -X GET "http://localhost:8080/api/professionals/7c065dd1-22b9-4bed-82e2-be973cb6ea47/appointments?status=pending" \
  -H "Content-Type: application/json"
```

### 4. Confirm Appointment
**PATCH** `/api/professionals/{id}/appointments/{appointment_id}/confirm`

Confirm a pending appointment.

**cURL:**
```bash
curl -X PATCH "http://localhost:8080/api/professionals/7c065dd1-22b9-4bed-82e2-be973cb6ea47/appointments/71a738d8-6695-4fa3-b68a-c58797801258/confirm" \
  -H "Content-Type: application/json"
```

### 5. Cancel Appointment (Professional)
**PATCH** `/api/professionals/{id}/appointments/{appointment_id}/cancel`

Cancel an appointment as a professional.

**Request Body:**
```json
{
  "cancellation_reason": "Client requested to reschedule"
}
```

**cURL:**
```bash
curl -X PATCH "http://localhost:8080/api/professionals/7c065dd1-22b9-4bed-82e2-be973cb6ea47/appointments/71a738d8-6695-4fa3-b68a-c58797801258/cancel" \
  -H "Content-Type: application/json" \
  -d '{
    "cancellation_reason": "Client requested to reschedule"
  }'
```

### 6. Create Unavailable Appointment
**POST** `/api/professionals/{id}/unavailable_appointments`

Mark a time period as unavailable.

**Request Body:**
```json
{
  "start_at": "2024-01-20T09:00:00Z",
  "end_at": "2024-01-20T17:00:00Z"
}
```

**cURL:**
```bash
curl -X POST "http://localhost:8080/api/professionals/7c065dd1-22b9-4bed-82e2-be973cb6ea47/unavailable_appointments" \
  -H "Content-Type: application/json" \
  -d '{
    "start_at": "2024-01-20T09:00:00Z",
    "end_at": "2024-01-20T17:00:00Z"
  }'
```

### 7. Get Professional Availability
**GET** `/api/professionals/{id}/availability`

Get hourly availability slots for a specific date (5:00-23:00).

**Query Parameters:**
- `date` (required): Date in YYYY-MM-DD format

**cURL:**
```bash
curl -X GET "http://localhost:8080/api/professionals/7c065dd1-22b9-4bed-82e2-be973cb6ea47/availability?date=2024-01-15" \
  -H "Content-Type: application/json"
```

**Response:**
```json
{
  "date": "2024-01-15",
  "slots": [
    {
      "start_time": "2024-01-15T05:00:00Z",
      "end_time": "2024-01-15T06:00:00Z",
      "available": true
    },
    {
      "start_time": "2024-01-15T10:00:00Z",
      "end_time": "2024-01-15T11:00:00Z",
      "available": false,
      "type": "appointment"
    },
    {
      "start_time": "2024-01-15T14:00:00Z",
      "end_time": "2024-01-15T15:00:00Z",
      "available": false,
      "type": "unavailable"
    }
  ]
}
```

## Appointment Management

### 1. Create Appointment
**POST** `/api/appointments`

Create a new appointment between a client and professional.

**Request Body:**
```json
{
  "client_id": "28c31a08-f740-440e-a161-6c8136478e2b",
  "professional_id": "7c065dd1-22b9-4bed-82e2-be973cb6ea47",
  "start_time": "2024-01-15T10:00:00Z",
  "end_time": "2024-01-15T11:00:00Z"
}
```

**cURL:**
```bash
curl -X POST "http://localhost:8080/api/appointments" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "28c31a08-f740-440e-a161-6c8136478e2b",
    "professional_id": "7c065dd1-22b9-4bed-82e2-be973cb6ea47",
    "start_time": "2024-01-15T10:00:00Z",
    "end_time": "2024-01-15T11:00:00Z"
  }'
```

## Database Schema

### Clients Table
```sql
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Professionals Table
```sql
CREATE TABLE professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    phone_number VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Appointments Table
```sql
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type appointment_type NOT NULL,           -- 'appointment' or 'unavailable'
    client_id UUID REFERENCES clients(id),   -- Can be null for unavailable
    professional_id UUID NOT NULL REFERENCES professionals(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status appointment_status DEFAULT 'pending',
    cancellation_reason TEXT,
    cancelled_by_professional_id UUID REFERENCES professionals(id),
    cancelled_by_client_id UUID REFERENCES clients(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Enums
```sql
CREATE TYPE appointment_type AS ENUM ('appointment', 'unavailable');
CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');
```

## Development Commands

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

## Accessing Services

- **API Server**: http://localhost:8080
- **pgAdmin**: http://localhost:8081
  - Email: admin@booking.com
  - Password: admin
- **PostgreSQL**: localhost:5432
  - Database: booking_db
  - User: booking_user
  - Password: booking_pass

## Key Features

### Appointment Status Flow
1. **Pending**: Newly created appointment (default)
2. **Confirmed**: Professional has confirmed the appointment
3. **Cancelled**: Either client or professional cancelled
4. **Completed**: Appointment has been completed

### Availability System
- **Hourly Slots**: 5:00 AM to 11:00 PM (18 slots per day)
- **Conflict Detection**: Automatically detects overlapping appointments
- **Type Detection**: Shows whether slot is blocked by appointment or unavailable period

### Cancellation System
- **Client Cancellation**: Clients can cancel their own appointments
- **Professional Cancellation**: Professionals can cancel any appointment
- **Reason Required**: Both must provide cancellation reason
- **Status Tracking**: Tracks who cancelled the appointment

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": {
    "code": "validation_error",
    "message": "Invalid request body",
    "details": "Field 'first_name' is required"
  }
}
```

Common error codes:
- `validation_error`: Invalid input data
- `database_error`: Database operation failed
- `not_found`: Resource not found
- `unauthorized`: Authentication required
- `conflict`: Resource already exists

## Security Considerations

- **Password Hashing**: Using bcrypt for professional passwords
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Type-safe queries with SQLC
- **UUID Usage**: All IDs are UUIDs for security

## Contributing

1. Follow the existing code structure and patterns
2. Add tests for new features
3. Update documentation for API changes
4. Use conventional commit messages

## License

[Add your license here]
