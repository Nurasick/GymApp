# GymApp Backend

Production-grade health & fitness backend API built with Go, Echo framework, and Clean Architecture.

## Features

- **User Management**: Registration, login, JWT authentication with refresh tokens
- **Recipe Generation**: AI-powered recipe recommendations from text or image ingredients
- **Training Plans**: Personalized workout plans based on user metrics
- **PostgreSQL**: Full database integration with migrations
- **Docker**: Complete containerized setup with docker-compose
- **Clean Architecture**: Domain, repository, service, and handler layers
- **Security**: bcrypt password hashing, JWT tokens, CORS protection

## Project Structure

```
GymApp/
├── cmd/api/
│   └── main.go                 # Application entry point
├── internal/
│   ├── domain/                 # Business entities and interfaces
│   │   ├── user.go
│   │   ├── recipe.go
│   │   ├── training.go
│   │   └── refresh_token.go
│   ├── repository/
│   │   └── postgres/           # Database implementations
│   │       ├── user_repo.go
│   │       ├── recipe_repo.go
│   │       ├── training_repo.go
│   │       └── refresh_token_repo.go
│   ├── service/                # Business logic
│   │   ├── auth_service.go
│   │   ├── recipe_service.go
│   │   ├── training_service.go
│   │   └── ai_service.go
│   ├── handler/
│   │   └── http/               # HTTP handlers and routes
│   │       ├── auth_handler.go
│   │       ├── recipe_handler.go
│   │       ├── training_handler.go
│   │       └── health_handler.go
│   ├── middleware/
│   │   └── auth.go             # JWT authentication middleware
│   ├── config/
│   │   └── config.go           # Configuration loader
│   └── database/
│       └── db.go               # Database connection pool
├── migrations/                 # Goose SQL migrations
├── pkg/utils/
│   └── logger.go               # Logging utilities
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── .env.example
```

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Echo v4
- **Database**: PostgreSQL 15
- **ORM/Query**: pgx/v5
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Migrations**: Goose
- **Container**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker Compose)

### Local Development

1. **Clone and setup**
```bash
cd c:\Users\Nurali\GoProjects\GymApp
cp .env.example .env
```

2. **Install dependencies**
```bash
go mod download
go mod tidy
```

3. **Start PostgreSQL**
```bash
docker-compose up -d postgres
```

The API will start on `http://localhost:8080`

### Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f api
```

## Configuration

Create a `.env` file from `.env.example`:

```env
SERVER_PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gymapp
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-change-in-production

AI_API_KEY=your-openai-key
AI_BASE_URL=https://api.openai.com/v1
AI_MODEL=gpt-3.5-turbo
```

## API Endpoints

See [API_DOCS.md](API_DOCS.md) for complete API documentation.

### Health Check
- `GET /health` - Returns server status

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user
- `POST /auth/refresh` - Refresh access token

### Recipes
- `POST /recipes/from-text` - Generate recipes from text ingredients
- `POST /recipes/from-image` - Generate recipes from image (multipart)
- `GET /recipes/history` - Get user's recipe history

### Training Plans
- `POST /training/generate` - Generate personalized training plan
- `GET /training/latest` - Get latest training plan

## Database Schema

### users
- id (BIGSERIAL PK)
- email (VARCHAR 255, UNIQUE)
- password_hash (TEXT)
- height (INT)
- weight (INT)
- goal (VARCHAR 100)
- created_at (BIGINT)
- updated_at (BIGINT)

### recipes
- id (BIGSERIAL PK)
- user_id (BIGINT FK → users)
- ingredients (TEXT)
- ai_response (TEXT)
- created_at (BIGINT)

### training_plans
- id (BIGSERIAL PK)
- user_id (BIGINT FK → users)
- plan_json (TEXT)
- created_at (BIGINT)

### refresh_tokens
- id (BIGSERIAL PK)
- user_id (BIGINT FK → users)
- token (TEXT, UNIQUE)
- expires_at (BIGINT)
- created_at (BIGINT)

## Security Considerations

1. **JWT Secret**: Change `JWT_SECRET` in production
2. **Database**: Use strong passwords and SSL connections
3. **API Key**: Secure your AI API key in environment variables
4. **CORS**: Configure allowed origins based on your frontend
5. **Rate Limiting**: Consider adding rate limiting middleware for production
6. **Input Validation**: All endpoints validate input data
7. **Password Hashing**: Using bcrypt with default cost

## Error Handling

All errors are returned in JSON format:

```json
{
  "message": "error description"
}
```

HTTP Status Codes:
- `200 OK` - Successful request
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing/invalid authentication
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Logging

The application uses structured logging with timestamps. Logs are output to stdout and include:
- INFO: General application flow
- ERROR: Error conditions
- DEBUG: Detailed debugging (development only)


