include .env

.PHONY: build-backend
build-backend:
	$(MAKE) -C backend ogen
	$(MAKE) -C backend sqlc
	$(MAKE) -C backend build

.PHONY: test-backend
test-backend:
	$(MAKE) -C backend ogen
	$(MAKE) -C backend sqlc
	$(MAKE) -C backend test

.PHONY: build-image
build-image:
	$(MAKE) -C backend build-image

.PHONY: ogen
ogen:
	$(MAKE) -C backend ogen
