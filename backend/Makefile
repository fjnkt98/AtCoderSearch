include .env

.PHONY: sqlc
sqlc:
	@atlas schema inspect -u ${DATABASE_URL} --format '{{ sql . }}' > schema.sql
	@sqlc generate

.PHONY: sqldef
sqldef:
	@psqldef --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} --port=${DB_PORT} --file=${DB_CONFIG} atcodersearch

.PHONY: atlas-schema-inspect
atlas-schema-inspect:
	@atlas schema inspect -u ${DATABASE_URL}

.PHONY: atlas-schema-apply
atlas-schema-apply:
	@atlas schema apply -u ${DATABASE_URL} --to file://schema.hcl
