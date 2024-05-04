include .env

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: sqldef
sqldef:
	@psqldef --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} --port=${DB_PORT} --file=${DB_CONFIG} atcodersearch
