include .env

.PHONY: build-backend
build-backend:
	$(MAKE) -C backend sqlc
	$(MAKE) -C backend build

.PHONY: test-backend
test-backend:
	$(MAKE) -C backend test

.PHONY: build-image
build-image:
	$(MAKE) -C backend build-image

.PHONY: buf-generate
buf-generate:
	$(MAKE) -C proto buf-generate

.PHONY: buf-lint
buf-lint:
	$(MAKE) -C proto buf-lint
