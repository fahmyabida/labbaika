.PHONY: migrate-create migrate-up migrate-down migrate-status migrate-force

MIGRATION_DIR := pkg/database/migration

# Host-side defaults. The container's .env uses PG_HOST=postgres (docker DNS),
# which does not resolve from the host, so migrations use localhost instead.
# Override per-invocation: make migrate-up PG_HOST=10.0.0.5 PG_DATABASE=staging
PG_HOST     ?= localhost
PG_PORT     ?= 5432
PG_USER     ?= postgres
PG_PASSWORD ?= postgres
PG_DATABASE ?= labbaika
PG_SSLMODE  ?= disable

DATABASE_URL := postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=$(PG_SSLMODE)

define require_migrate
@command -v migrate >/dev/null || { echo "golang-migrate not installed. Run: brew install golang-migrate"; exit 1; }
endef

migrate-create:
	$(require_migrate)
	@test -n "$(name)" || { echo "Usage: make migrate-create name=create_users_table"; exit 1; }
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-up:
	$(require_migrate)
	@echo "Applying migrations to $(PG_DATABASE) @ $(PG_HOST):$(PG_PORT)"
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" up $(n)

migrate-down:
	$(require_migrate)
	@test -n "$(n)" || { \
		echo "Refusing to roll back every migration."; \
		echo "Usage: make migrate-down n=1        # roll back 1 step"; \
		echo "       make migrate-down n=all      # roll back everything (destructive)"; \
		exit 1; }
	@echo "Rolling back [$(n)] on $(PG_DATABASE) @ $(PG_HOST):$(PG_PORT)"
	@if [ "$(n)" = "all" ]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" down; \
	else \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" down $(n); \
	fi

migrate-status:
	$(require_migrate)
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" version

migrate-force:
	$(require_migrate)
	@test -n "$(v)" || { echo "Usage: make migrate-force v=1   # mark version v as clean after a failed migration"; exit 1; }
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" force $(v)
