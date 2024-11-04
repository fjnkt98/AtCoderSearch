.PHONY: test
test: ogen
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: ogen
ogen:
	ogen -package api -target api -clean api/openapi.yaml

.PHONY: build
build: sqlc ogen
	go build -o build/atcodersearch

VERSION:=$(shell cat VERSION)

.PHONY: build-image
build-image:
	docker build -t atcodersearch-backend:$(VERSION) .
