.PHONY: codegen
codegen:
	npx openapi-typescript ../backend/api/openapi.yaml -o ./app/types.d.ts

.PHONY: build
build: codegen
	npm run build

VERSION:=$(shell cat package.json | jq -r '.version')

.PHONY: build-image
build-image:
	docker build -t atcodersearch-frontend:$(VERSION) .
