DB_URL := "postgres://postgres:postgres@localhost:5432/gator"

migrate-up:
	goose -dir sql/schema postgres $(DB_URL) up

migrate-down:
	goose -dir sql/schema postgres $(DB_URL) down
