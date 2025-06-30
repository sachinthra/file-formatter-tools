# Variables
FRONTEND_DIR=frontend
BACKEND_DIR=backend
DEPLOY_DIR=deploy
DOCKER_COMPOSE=docker-compose -f $(DEPLOY_DIR)/docker-compose.yml

# Targets
.PHONY: all build-frontend build-backend docker-up docker-down clean

# Build both frontend and backend
all: build-frontend build-backend

# Build the frontend
build-frontend:
	@echo "Building the frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build

# Build the backend (if needed)
build-backend:
	@echo "Building the backend..."
	cd $(BACKEND_DIR) && go build -o main .

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	cd $(DEPLOY_DIR) && cp .env.example .env && $(DOCKER_COMPOSE) up --build

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	cd $(DEPLOY_DIR) && $(DOCKER_COMPOSE) down

# Clean up build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(FRONTEND_DIR)/dist
	rm -f $(BACKEND_DIR)/main