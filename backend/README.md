# Quizora Backend API

A robust REST API for the Quizora MCQ platform built with Go, Gin, and MySQL.

## Features

- 🔐 **JWT Authentication** - Secure token-based authentication
- 📱 **MSISDN Validation** - Bangladesh mobile number format validation
- 🔒 **Password Hashing** - Bcrypt password encryption
- 🛡️ **Input Validation** - Comprehensive request validation
- 🚀 **Clean Architecture** - Well-structured, maintainable codebase
- 📊 **Database Migration** - Automatic schema migration with GORM
- 🌐 **CORS Support** - Configured for frontend integration

## Tech Stack

- **Go 1.21+**
- **Gin Web Framework**
- **GORM** - Object Relational Mapping
- **MySQL 8.4+**
- **JWT** - JSON Web Tokens
- **Bcrypt** - Password hashing

## Project Structure

```
backend/
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migration
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Data models and DTOs
│   ├── services/        # Business logic layer
│   └── utils/           # Utility functions
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── .env.example         # Environment variables template
└── README.md           # This file
```

## Prerequisites

- Go 1.21 or higher
- MySQL 8.4+ running in Docker
- Git

## Setup Instructions

### 1. Clone and Navigate

```bash
cd /Users/bs01562/Documents/Personal/quizora/backend
```

### 2. Environment Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit the .env file with your MySQL credentials
# The default settings should work with your Docker MySQL setup
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Database Setup

Make sure your MySQL Docker container is running:

```bash
# Check if MySQL container is running
docker ps | grep mysql

# If not running, start it
docker start common-mysql-1
```

The API will automatically create the `quizora` database and run migrations on startup.

### 5. Run the Application

```bash
# Development mode
go run main.go

# Or build and run
go build -o quizora-backend
./quizora-backend
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

- `GET /health` - API health status

### Authentication

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/profile` - Get user profile (Protected)

## API Documentation

### Register User

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "msisdn": "+8801712345678",
  "password": "securepassword123"
}
```

**Response:**

```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "msisdn": "+8801712345678",
    "is_active": true,
    "created_at": "2025-01-26T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Login User

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "msisdn": "+8801712345678",
  "password": "securepassword123"
}
```

### Get Profile

```http
GET /api/v1/auth/profile
Authorization: Bearer <jwt_token>
```

## MSISDN Format Support

The API supports multiple Bangladesh mobile number formats:

- `+8801XXXXXXXXX` (preferred format)
- `8801XXXXXXXXX`
- `01XXXXXXXXX`

All formats are automatically normalized to `+8801XXXXXXXXX`.

## Security Features

- **Password Hashing**: Bcrypt with default cost
- **JWT Tokens**: 24-hour expiration
- **Input Validation**: Comprehensive validation for all inputs
- **CORS Protection**: Configured for specific frontend origin
- **SQL Injection Protection**: GORM provides built-in protection

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  name varchar(100) NOT NULL,
  msisdn varchar(20) NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
  is_active boolean DEFAULT true,
  created_at datetime(3),
  updated_at datetime(3),
  deleted_at datetime(3),
  INDEX idx_users_deleted_at (deleted_at)
);
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
# Build binary
go build -ldflags="-w -s" -o quizora-backend main.go

# Set production environment
export GIN_MODE=release
export JWT_SECRET=your-super-secure-production-secret
```

## Environment Variables

| Variable       | Description              | Default                     |
| -------------- | ------------------------ | --------------------------- |
| `DB_HOST`      | MySQL host               | `localhost`                 |
| `DB_PORT`      | MySQL port               | `3306`                      |
| `DB_USER`      | MySQL username           | `root`                      |
| `DB_PASSWORD`  | MySQL password           | `root`                      |
| `DB_NAME`      | Database name            | `quizora`                   |
| `JWT_SECRET`   | JWT signing secret       | `your-super-secret-jwt-key` |
| `PORT`         | Server port              | `8080`                      |
| `GIN_MODE`     | Gin mode (debug/release) | `debug`                     |
| `FRONTEND_URL` | Frontend URL for CORS    | `http://localhost:3000`     |

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error Type",
  "message": "Detailed error message"
}
```

Common HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `409` - Conflict (duplicate data)
- `500` - Internal Server Error

## Next Steps

1. **Install dependencies**: `go mod tidy`
2. **Start the server**: `go run main.go`
3. **Test with your frontend**: Update frontend API calls to `http://localhost:8080`
4. **Add more features**: Extend the API with quiz functionality

## Contributing

1. Follow Go conventions and best practices
2. Add tests for new features
3. Update documentation
4. Use meaningful commit messages

## License

This project is part of the Quizora platform.
