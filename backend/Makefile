.PHONY: build
build:
	@go build -o build/atcodersearch

.PHONY: test
test:
	@go test -cover ./... -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: build-image
build-image:
	docker build -t atcodersearch-batch:latest .

.PHONY: ogen
ogen:
	ogen -package api -target api -clean api/openapi.yaml
