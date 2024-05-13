.PHONY: migrate-up
migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -path db/migrations --verbose up

.PHONY: migrate-down
migrate-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -path db/migrations --verbose down