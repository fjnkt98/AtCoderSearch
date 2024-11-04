include .env

.PHONY: build-frontend
build-frontend:
	$(MAKE) -C frontend build

.PHONY: build-backend
build-backend:
	$(MAKE) -C backend build

.PHONY: build
build: build-backend build-frontend

.PHONY: test-backend
test-backend: build-backend
	$(MAKE) -C backend test

.PHONY: test
test: test-backend

.PHONY: build-image
build-image: build-backend
	$(MAKE) -C backend build-image
