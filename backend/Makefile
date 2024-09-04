include .env

.PHONY: test
test:
	@go test -cover ./... -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

.PHONY: sqlc
sqlc:
	@atlas schema inspect -u ${DATABASE_URL} --format '{{ sql . }}' > schema.sql
	@sqlc generate
