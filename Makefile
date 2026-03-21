.PHONY: help \
        backend-run backend-build backend-test backend-tidy backend-swagger backend-clean \
        frontend-install frontend-dev frontend-build frontend-preview \
        dev

BACKEND_DIR  := backend
FRONTEND_DIR := frontend
BINARY       := $(BACKEND_DIR)/smart-allocation

# ─── Help ────────────────────────────────────────────────────────────────────

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Backend"
	@echo "  backend-run      Run the Go server (loads .env)"
	@echo "  backend-build    Compile the Go binary"
	@echo "  backend-test     Run all Go tests"
	@echo "  backend-tidy     Tidy Go module dependencies"
	@echo "  backend-swagger  Regenerate Swagger docs"
	@echo "  backend-clean    Remove compiled binary"
	@echo ""
	@echo "Frontend"
	@echo "  frontend-install Install npm dependencies"
	@echo "  frontend-dev     Start Vite dev server"
	@echo "  frontend-build   Build for production"
	@echo "  frontend-preview Preview the production build"
	@echo ""
	@echo "Combined"
	@echo "  dev              Run backend and frontend dev servers concurrently"

# ─── Backend ─────────────────────────────────────────────────────────────────

backend-run:
	cd $(BACKEND_DIR) && go run main.go

backend-build:
	cd $(BACKEND_DIR) && go build -o smart-allocation .

backend-test:
	cd $(BACKEND_DIR) && go test ./...

backend-tidy:
	cd $(BACKEND_DIR) && go mod tidy

backend-swagger:
	cd $(BACKEND_DIR) && swag init -g main.go --output docs

backend-clean:
	rm -f $(BINARY)

# ─── Frontend ────────────────────────────────────────────────────────────────

frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-preview:
	cd $(FRONTEND_DIR) && npm run preview

# ─── Combined ────────────────────────────────────────────────────────────────

dev:
	@trap 'kill 0' SIGINT; \
	(cd $(BACKEND_DIR) && go run main.go) & \
	(cd $(FRONTEND_DIR) && npm run dev) & \
	wait
