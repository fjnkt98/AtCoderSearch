include .env

.PHONY: build-backend-batch
build-backend-batch:
	$(MAKE) -C backend/cmd/batch build

.PHONY: build-backend-server
build-backend-server:
	$(MAKE) -C backend/cmd/server build

.PHONY: build-backend
build-backend:
	$(MAKE) -C backend sqlc
	$(MAKE) -C backend/cmd/batch build
	$(MAKE) -C backend/cmd/server build

.PHONY: test-backend
test-backend:
	$(MAKE) -C backend test

.PHONY: sqlc
sqlc:
	$(MAKE) -C backend sqlc

.PHONY: build-image
build-image:
	$(MAKE) -C batch build-image
