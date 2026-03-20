# Smart Allocation

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

Application for managing and rebalancing a variable income investment portfolio. It allows registering assets with their desired allocation parameters and automatically calculates how much to invest (or reduce) in each position so the portfolio reaches the defined target percentages.

## Technology Stack

- **Language**: Go 1.22+

## Repository Structure

```
smart-allocation/
├── backend/                # REST API in Go
│   ├── main.go
│   ├── go.mod
│   ├── docs/               # Auto-generated Swagger documentation
│   └── internal/
│       ├── domain/         # Entities, repository interfaces, custom errors
│       ├── application/    # Use cases and DTOs
│       ├── infrastructure/ # HTTP handlers, SQLite repository, routes
│       └── configuration/  # Database and Swagger setup
└── frontend/               # Web interface (coming soon)
```

## Quick Start

### Installing Go

1. Download Go from the [official website](https://golang.org/dl/)
   - For macOS: Download the `.pkg` file and follow the installation wizard

2. Verify installation:
   ```bash
   go version
   ```

3. Set up your Go workspace:
   ```bash
   mkdir -p $HOME/go/{bin,src,pkg}
   ```

4. Add the following to your shell profile (`.bashrc`, `.zshrc`, etc.):
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```

## Project Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd smart-allocation
   ```

2. **Install dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

The server will be available at `http://localhost:8080`.

Data is automatically persisted in the `portfolio.db` file created inside `backend/`.

## API Documentation

Once the application is running, you can access the interactive API documentation:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`

### Updating

To update the Swagger documentation after making changes to API endpoints:

```bash
cd backend
swag init -g main.go --output docs
```

## Development

### Available Commands

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Update Swagger documentation
swag init -g main.go --output docs
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
