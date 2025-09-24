# 🎮 Pokemon API

A robust REST API for managing Pokemon data with external PokeAPI integration, built in Go using hexagonal architecture and SOLID principles.

> **🚀 Key Features:** Complete REST API, PokeAPI integration, hexagonal architecture, automatic Swagger documentation, and Docker containerization.

## 📋 What's Inside

This project demonstrates modern software development practices:

- **🏗️ Hexagonal Architecture**: Clean separation between business logic and external concerns
- **🔌 SOLID Principles**: Dependency inversion, single responsibility, and interface segregation
- **🧪 Testable Design**: Mock-friendly interfaces for comprehensive testing
- **📚 Auto-Documentation**: Swagger UI generated from code annotations
- **🐳 Production Ready**: Docker containerization with PostgreSQL database
- **🔄 External Integration**: Real-time data fetching from PokeAPI

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd amaris-pokemon-challenge

# Start the application with Docker Compose
docker compose up -d

# The API will be available at:
# - API: http://localhost:8080
# - Swagger Documentation: http://localhost:8080/swagger/index.html
# - Health Check: http://localhost:8080/health
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Start PostgreSQL (you need to have it installed)
# Update connection details in main.go if needed

# Run the application
go run cmd/api/main.go
```

## 📚 API Endpoints

### Create Pokemon
```bash
curl -X POST http://localhost:8080/api/v1/pokemon \
  -H "Content-Type: application/json" \
  -d '{
    "name": "pikachu",
    "type1": "electric",
    "type2": ""
  }'
```

### Create Pokemon (Flexible Format)
```bash
# Format 1: Direct name
curl -X POST http://localhost:8080/api/v1/pokemon \
  -H "Content-Type: application/json" \
  -d '{
    "name": "charizard",
    "type1": "fire",
    "type2": "flying"
  }'

# Format 2: Nested pokemon object
curl -X POST http://localhost:8080/api/v1/pokemon \
  -H "Content-Type: application/json" \
  -d '{
    "pokemon": {
      "name": "blastoise"
    },
    "type1": "water",
    "type2": ""
  }'
```

### Get Pokemon by ID
```bash
curl http://localhost:8080/api/v1/pokemon/1
```

### List All Pokemon
```bash
curl http://localhost:8080/api/v1/pokemon
```

### Health Check
```bash
curl http://localhost:8080/health
```

## 🏗️ Architecture

This project implements **Hexagonal Architecture** with the following structure:

### 📊 Architecture Diagram

![Hexagonal Architecture](docs/plantuml/Pokemon%20API%20-%20Hexagonal%20Architecture.png)

*Hexagonal architecture separates business logic (Core) from external concerns (Adapters), enabling maximum flexibility and testability.*

```
pokemon-api/
├── cmd/api/                    # Application entry point
├── internal/
│   ├── core/                   # Business logic (inner layer)
│   │   ├── domain/            # Entities
│   │   ├── ports/             # Interfaces
│   │   └── services/          # Use cases
│   └── adapters/              # External adapters (outer layer)
│       ├── handlers/          # HTTP handlers
│       ├── repositories/      # Database
│       └── external/          # External API clients
├── docs/                      # Swagger documentation
├── docker-compose.yml
├── Dockerfile
└── README.md
```

### 🔄 Service Flow

![Service Flow](docs/plantuml/Pokemon%20Service%20Flow.png)

*This diagram shows how data flows from HTTP requests to the database, passing through business logic and external integrations.*

### 🎯 Why Hexagonal Architecture?

**Hexagonal architecture** (also called "Ports and Adapters") was chosen because:

- **🔒 Isolation**: Business logic is completely separated from technical concerns
- **🧪 Testability**: Easy to create mocks and unit tests
- **🔧 Maintainability**: Changes to DB or external APIs don't affect business core
- **📈 Scalability**: Easy to add new adapters without modifying existing code
- **🏗️ SOLID**: Implements all SOLID programming principles

### 📁 Layer Structure

```
🧠 Core (Business Logic)
├── 📋 Domain: Entities and business rules
├── 🔌 Ports: Interfaces (contracts)
└── ⚙️ Services: Use cases and application logic

🔌 Adapters (Implementations)
├── 🌐 HTTP Handlers: Web request handling
├── 🗄️ Repositories: Database access
└── 🔗 External Clients: External API integration
```

## 🛠️ Technology Stack

- **Language**: Go 1.21
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **External API**: PokeAPI (https://pokeapi.co/)
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker + Docker Compose
- **Architecture**: Hexagonal Architecture
- **Principles**: SOLID principles

## ✨ Key Features

### 🎮 **API Functionality**
- ✅ **Complete REST API** with Gin framework
- ✅ **PokeAPI integration** to fetch official Pokemon data
- ✅ **Flexible name handling** (direct or nested format)
- ✅ **Robust input validation** for data integrity
- ✅ **Automatic documentation** with Swagger UI

### 🏗️ **Architecture & Quality**
- ✅ **Hexagonal architecture** properly implemented
- ✅ **SOLID principles** applied throughout the codebase
- ✅ **Clear separation** between business logic and adapters
- ✅ **Well-defined interfaces** for maximum flexibility

### 🛠️ **Infrastructure & DevOps**
- ✅ **PostgreSQL database** with GORM ORM
- ✅ **Complete Docker containerization**
- ✅ **Docker Compose** for local development
- ✅ **Health check endpoint** for monitoring
- ✅ **Error handling** with appropriate HTTP status codes

## 🔧 Configuration

The application uses environment variables for configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host |
| `DB_USER` | `pokemon_user` | Database user |
| `DB_PASSWORD` | `pokemon_pass` | Database password |
| `DB_NAME` | `pokemon_db` | Database name |
| `DB_PORT` | `5432` | Database port |
| `POKEAPI_BASE_URL` | `https://pokeapi.co/api/v2` | PokeAPI base URL |
| `PORT` | `8080` | Application port |

## 🧪 Testing

```bash
# Run tests (when implemented)
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 📖 API Documentation

Visit http://localhost:8080/swagger/index.html for interactive API documentation.

## 🚨 Important Notes

1. **Pokemon Name Handling**: The API supports both direct name format and nested `pokemon.name` format for maximum flexibility.

2. **External API**: The application fetches Pokemon data from PokeAPI and stores it locally in the database.

3. **Database**: Uses PostgreSQL with GORM for data persistence and automatic migrations.

4. **Error Handling**: Returns appropriate HTTP status codes (400, 404, 409, 500) with descriptive error messages.

## 🐳 Docker Commands

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild and start
docker-compose up -d --build
```

## 🔍 Development

### Local Development Setup

1. Install Go 1.21+
2. Install PostgreSQL
3. Clone the repository
4. Run `go mod download`
5. Update database connection in `main.go` if needed
6. Run `go run cmd/api/main.go`

### Code Structure

- **Domain**: Core business entities and interfaces
- **Services**: Business logic and use cases
- **Handlers**: HTTP request/response handling
- **Repositories**: Data access layer
- **External**: Third-party API integrations

## 📄 License

This project is part of the Amaris Pokemon Challenge.

---

**Built with ❤️ using Go and Hexagonal Architecture**
