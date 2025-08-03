GOOSE=goose
DB_URL=postgres://simt_user:simtDev@2025!@localhost:5432/simt_dev?sslmode=disable
TEST_DB_URL=postgres://simt_user:simtDev@2025!@localhost:5432/simt_test?sslmode=disable
MIGRATIONS_DIR=./migrations

.PHONY: migrate-up migrate-down migrate-status create-migration

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

create-migration:
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	filename="$${timestamp}_$${name}.sql"; \
	touch $(MIGRATIONS_DIR)/$$filename; \
	echo "-- +goose Up\n\n-- +goose Down" > $(MIGRATIONS_DIR)/$$filename; \
	echo "Created $$filename"

