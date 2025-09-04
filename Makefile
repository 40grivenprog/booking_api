# Simple commands
.PHONY: start
start: ## Start the app
	docker compose up

.PHONY: stop
stop: ## Stop the app
	docker compose down

.PHONY: logs
logs: ## Show logs
	docker compose logs -f

.PHONY: build
build: ## Build the app
	docker compose build

.PHONY: restart
restart: ## Restart the app
	docker compose restart

.PHONY: clean
clean: ## Clean everything
	docker compose down -v
	docker system prune -f

# PostgreSQL commands
.PHONY: pgadmin
pgadmin: ## Open pgAdmin in browser
	@echo "Opening pgAdmin at http://localhost:8081"
	@if command -v open >/dev/null 2>&1; then open http://localhost:8081\; fi

.PHONY: db-shell
db-shell: ## Connect to PostgreSQL shell
	docker compose exec postgres psql -U booking_user -d booking_db

.PHONY: status
status: ## Show status of all services
	docker compose ps

# Migration commands
.PHONY: migrate-build
migrate-build: ## Build migration binary
	GOOS=linux GOARCH=amd64 go build -o bin/migrate ./cmd/migrate

.PHONY: migrate-up
migrate-up: migrate-build ## Run database migrations up
	docker run --rm --network booking_api_default -v $(PWD)/bin:/bin -v $(PWD)/internal/migrations:/migrations -e DATABASE_URL="postgres://booking_user:booking_pass@postgres:5432/booking_db?sslmode=disable" alpine:latest /bin/migrate -direction=up

.PHONY: migrate-down
migrate-down: migrate-build ## Run database migrations down
	docker run --rm --network booking_api_default -v $(PWD)/bin:/bin -v $(PWD)/internal/migrations:/migrations -e DATABASE_URL="postgres://booking_user:booking_pass@postgres:5432/booking_db?sslmode=disable" alpine:latest /bin/migrate -direction=down

.PHONY: migrate-create
migrate-create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./internal/migrations -seq $$name

# SQLC commands
.PHONY: sqlc-generate
sqlc-generate: ## Generate SQLC code
	sqlc generate

.PHONY: sqlc-validate
sqlc-validate: ## Validate SQLC configuration
	sqlc compile

# Development commands
.PHONY: dev
dev: ## Start development environment
	docker compose -f docker-compose.dev.yml up

.PHONY: dev-build
dev-build: ## Build development environment
	docker compose -f docker-compose.dev.yml build

.PHONY: dev-stop
dev-stop: ## Stop development environment
	docker compose -f docker-compose.dev.yml down
