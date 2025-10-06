# Define the database variables
DB_NAME=merchant_service
DB_USER=user
DB_PASS=password
DB_HOST=localhost
DB_PORT=5433

# Target for creating database
create-db:
	@echo "Creating database if not exists..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d postgres -c 'CREATE DATABASE "$(DB_NAME)" WITH OWNER "$(DB_USER)" TEMPLATE template0 ENCODING UTF8;'

# Target for creating new migrations
migrate-create:
	@migrate create -ext sql -dir scripts/migrations -seq $(name)

# Target for applying migrations
migrate-up:
	@echo "Applying migrations..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/migrations/000001_create-merchant-table.up.sql
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/migrations/000003_create-item-table.up.sql

# Target for reverting migrations
migrate-down:
	@echo "Reverting migrations..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/migrations/000001_create-merchant-table.down.sql
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/migrations/000003_create-item-table.down.sql

# Target for dropping database (force disconnect users first)
drop-db:
	@echo "Dropping database..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$(DB_NAME)' AND pid <> pg_backend_pid();"
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d postgres -c 'DROP DATABASE IF EXISTS "$(DB_NAME)";'

# Target for checking current state
check-db:
	@echo "Checking database and tables..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"

# Target for manual table drop (if needed)
drop-table:
	@echo "Dropping users table..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) -c "DROP TABLE IF EXISTS merchant;"

# Clear all Table Data
clear-merchant-table:
	@echo "Clearing all data from users table..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) -c "DELETE FROM merchants;"

# Target for resetting database
reset-db: drop-db create-db migrate-up
	@echo "Database reset complete"

# Target for generating sqlc code
sqlc-generate:
	@echo "Generating sqlc code..."
	@sqlc generate

# seeding DB
seed-db-up:
	@echo "Seeding database with initial data..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/migrations/000002_seed-merchant-data.up.sql

seed-db-down:
	@echo "Seeding database with initial data..."
	@docker exec -i belimang-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/migrations/000002_seed-merchant-data.down.sql

.PHONY: create-db migrate-create migrate-up migrate-down drop-db check-db drop-table reset-db sqlc-generate