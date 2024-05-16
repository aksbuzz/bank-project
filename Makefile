.PHONY: migrate-up
migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/library?sslmode=disable" -path migration --verbose up

.PHONY: migrate-down
migrate-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/library?sslmode=disable" -path migration --verbose down