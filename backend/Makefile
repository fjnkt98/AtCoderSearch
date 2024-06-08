include .env

.PHONY: sqlc
sqlc:
	@atlas schema inspect -u ${DATABASE_URL} --format '{{ sql . }}' > schema.sql
	@sqlc generate

.PHONY: atlas-schema-inspect
atlas-schema-inspect:
	@atlas schema inspect -u ${DATABASE_URL}

.PHONY: atlas-schema-apply
atlas-schema-apply:
	@atlas schema apply -u ${DATABASE_URL} --to file://schema.hcl
