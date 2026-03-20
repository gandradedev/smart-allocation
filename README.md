# Smart Allocation

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
![technology React](https://img.shields.io/badge/technology-react-61DAFB.svg)

Application for managing and rebalancing a variable income investment portfolio. It allows registering assets with their desired allocation parameters and automatically calculates how much to invest (or reduce) in each position so the portfolio reaches the defined target percentages.

## Technology Stack

- **Backend**: Go 1.22+
- **Frontend**: React 18 + TypeScript + Vite
- **Styling**: Tailwind CSS
- **State/Data fetching**: TanStack Query
- **Form validation**: React Hook Form + Zod

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
└── frontend/               # React web application
    ├── index.html
    ├── vite.config.ts      # Dev server + proxy to backend
    └── src/
        ├── types/          # TypeScript interfaces
        ├── services/       # API calls
        ├── hooks/          # TanStack Query hooks
        ├── utils/          # Number formatting (BRL)
        └── components/     # UI components
```

## Quick Start

### Prerequisites

- [Go 1.22+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)

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

2. **Configure the brapi.dev token**

   The backend requires a token to fetch asset prices from [brapi.dev](https://brapi.dev).

   Copy the example file and set your token:
   ```bash
   cd backend
   cp .env.example .env
   ```

   Then edit `backend/.env` and replace the placeholder with your token:
   ```
   BRAPI_TOKEN=your_token_here
   ```

   > Get your free token at [brapi.dev](https://brapi.dev) after creating an account.
   > The `.env` file is ignored by git and should never be committed.

3. **Start the backend**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

   The API will be available at `http://localhost:8080`.
   Data is automatically persisted in `backend/portfolio.db`.

3. **Start the frontend** (in a separate terminal)
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

   The app will be available at `http://localhost:5173`.

## API Documentation

Once the backend is running, the interactive API documentation is available at:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`

### Updating Swagger

To regenerate the documentation after making changes to API endpoints:

```bash
cd backend
swag init -g main.go --output docs
```

## Development

### Backend

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Update Swagger documentation
swag init -g main.go --output docs
```

### Frontend

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build
```

## Testing

### Backend

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
