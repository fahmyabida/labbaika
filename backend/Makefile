.PHONY: migrate create

# Phony target that allows for "make migrate create"
migrate: create

# Create migration target
create:
	@echo "Creating migration with name: $(name)"
	@migrate create -ext sql -dir /pkg/database/migration/ -seq $(name)
	# @go run pkg/database/migration/migrate.go create -ext sql -dir pkg/database/migration/ -seq $(name)
