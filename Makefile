migrate-up:
	docker compose --profile tools run --rm migrate up

migrate-down:
	docker compose --profile tools run --rm migrate down 1

migrate-version:
	docker compose --profile tools run --rm migrate version