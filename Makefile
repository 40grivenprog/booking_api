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

.PHONY: migrate-version
migrate-version: migrate-build ## Check migration version
	docker run --rm --network booking_api_default -v $(PWD)/bin:/bin -v $(PWD)/internal/migrations:/migrations -e DATABASE_URL="postgres://booking_user:booking_pass@postgres:5432/booking_db?sslmode=disable" alpine:latest /bin/migrate -direction=up -dry-run

.PHONY: migrate-force
migrate-force: migrate-build ## Force migration version (use with caution)
	@read -p "Enter version to force: " version; \
	docker run --rm --network booking_api_default -v $(PWD)/bin:/bin -v $(PWD)/internal/migrations:/migrations -e DATABASE_URL="postgres://booking_user:booking_pass@postgres:5432/booking_db?sslmode=disable" alpine:latest /bin/migrate -direction=up -version=$$version -force

# Production migration commands
.PHONY: migrate-prod-up
migrate-prod-up: ## Run migrations up in production
	@echo "Running migrations for prod environment..."
	@cd ../booking_deployments/terraform/environments/prod && \
	NETWORK_INFO=$$(terraform output -json network_info) && \
	SUBNET_IDS=$$(echo $$NETWORK_INFO | jq -r '.subnet_ids | join(",")') && \
	SECURITY_GROUP_ID=$$(echo $$NETWORK_INFO | jq -r '.security_group_id') && \
	aws ecs run-task \
		--cluster "booking-app-prod" \
		--task-definition "booking-app-prod-db-migration" \
		--launch-type "FARGATE" \
		--network-configuration "awsvpcConfiguration={subnets=[$$SUBNET_IDS],securityGroups=[$$SECURITY_GROUP_ID],assignPublicIp=ENABLED}" \
		--overrides 'containerOverrides=[{name="db-migration",command=["-direction=up"]}]' \
		--region "eu-central-1"

.PHONY: migrate-prod-down
migrate-prod-down: ## Run migrations down in production
	@read -p "How many steps to rollback? (default: 1): " steps; \
	echo "Running migrations down for prod environment..."; \
	cd ../booking_deployments/terraform/environments/prod && \
	NETWORK_INFO=$$(terraform output -json network_info) && \
	SUBNET_IDS=$$(echo $$NETWORK_INFO | jq -r '.subnet_ids | join(",")') && \
	SECURITY_GROUP_ID=$$(echo $$NETWORK_INFO | jq -r '.security_group_id') && \
	aws ecs run-task \
		--cluster "booking-app-prod" \
		--task-definition "booking-app-prod-db-migration" \
		--launch-type "FARGATE" \
		--network-configuration "awsvpcConfiguration={subnets=[$$SUBNET_IDS],securityGroups=[$$SECURITY_GROUP_ID],assignPublicIp=ENABLED}" \
		--overrides "containerOverrides=[{name=\"db-migration\",command=[\"-direction=down\",\"-steps=$${steps:-1}\"]}]" \
		--region "eu-central-1"

.PHONY: migrate-prod-force
migrate-prod-force: ## Force migration version in production
	@read -p "Enter version to force: " version; \
	echo "Forcing migration version $$version for prod environment..."; \
	cd ../booking_deployments/terraform/environments/prod && \
	NETWORK_INFO=$$(terraform output -json network_info) && \
	SUBNET_IDS=$$(echo $$NETWORK_INFO | jq -r '.subnet_ids | join(",")') && \
	SECURITY_GROUP_ID=$$(echo $$NETWORK_INFO | jq -r '.security_group_id') && \
	aws ecs run-task \
		--cluster "booking-app-prod" \
		--task-definition "booking-app-prod-db-migration" \
		--launch-type "FARGATE" \
		--network-configuration "awsvpcConfiguration={subnets=[$$SUBNET_IDS],securityGroups=[$$SECURITY_GROUP_ID],assignPublicIp=ENABLED}" \
		--overrides "containerOverrides=[{name=\"db-migration\",command=[\"-direction=force\",\"-version=$$version\"]}]" \
		--region "eu-central-1"

.PHONY: migrate-dev-up
migrate-dev-up: ## Run migrations up in development
	@echo "Running migrations for dev environment..."
	@cd ../booking_deployments/terraform/environments/dev && \
	NETWORK_INFO=$$(terraform output -json network_info) && \
	SUBNET_IDS=$$(echo $$NETWORK_INFO | jq -r '.subnet_ids | join(",")') && \
	SECURITY_GROUP_ID=$$(echo $$NETWORK_INFO | jq -r '.security_group_id') && \
	aws ecs run-task \
		--cluster "booking-app-dev" \
		--task-definition "booking-app-dev-db-migration" \
		--launch-type "FARGATE" \
		--network-configuration "awsvpcConfiguration={subnets=[$$SUBNET_IDS],securityGroups=[$$SECURITY_GROUP_ID],assignPublicIp=ENABLED}" \
		--overrides 'containerOverrides=[{name="db-migration",command=["-direction=up"]}]' \
		--region "eu-central-1"

.PHONY: migrate-dev-down
migrate-dev-down: ## Run migrations down in development
	@read -p "How many steps to rollback? (default: 1): " steps; \
	echo "Running migrations down for dev environment..."; \
	cd ../booking_deployments/terraform/environments/dev && \
	NETWORK_INFO=$$(terraform output -json network_info) && \
	SUBNET_IDS=$$(echo $$NETWORK_INFO | jq -r '.subnet_ids | join(",")') && \
	SECURITY_GROUP_ID=$$(echo $$NETWORK_INFO | jq -r '.security_group_id') && \
	aws ecs run-task \
		--cluster "booking-app-dev" \
		--task-definition "booking-app-dev-db-migration" \
		--launch-type "FARGATE" \
		--network-configuration "awsvpcConfiguration={subnets=[$$SUBNET_IDS],securityGroups=[$$SECURITY_GROUP_ID],assignPublicIp=ENABLED}" \
		--overrides "containerOverrides=[{name=\"db-migration\",command=[\"-direction=down\",\"-steps=$${steps:-1}\"]}]" \
		--region "eu-central-1"

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

# Database commands
.PHONY: db-exec
db-exec: ## Connect to database via ECS Exec
	@echo "Connecting to database via ECS Exec..."
	@echo "Getting running task..."
	TASK_ARN=$$(aws ecs list-tasks --cluster booking-app-dev --service-name booking-api-dev --region eu-central-1 --query 'taskArns[0]' --output text) && \
	if [ "$$TASK_ARN" = "None" ] || [ -z "$$TASK_ARN" ]; then \
			echo "No running tasks found. Make sure the service is running."; \
		exit 1; \
	fi && \
	echo "Task: $$TASK_ARN" && \
	echo "Connecting to database..." && \
	aws ecs execute-command \
		--cluster booking-app-dev \
		--task $$TASK_ARN \
		--container booking-api \
		--interactive \
		--command "sh -c 'psql -h \$$DB_HOST -U \$$DB_USER -d \$$DB_NAME'" \
		--region eu-central-1

