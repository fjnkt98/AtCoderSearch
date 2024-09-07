.PHONY: test
test:
	@go test -cover ./... -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

.PHONY: sqlc
sqlc:
	@sqlc generate
