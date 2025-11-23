.PHONY: run up down clean migration goose-up dt0


run:
	cd ./avito_test && godotenv -f ../.env go run ./cmd/avito_test/main.go

up:
	docker-compose up -d

down:
	docker-compose down

clean:
	docker-compose down -v
	docker system prune -f

migration:
	$(eval name := $(word 2, $(MAKECMDGOALS)))
	@if [ -z "$(name)" ]; then \
		echo "migration name is needed"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
		goose -dir ./migrations create $(name) sql 

goose-up:
	goose -dir ./migrations postgres "postgresql://test:test@localhost:5432/test?sslmode=disable" up

goose-down:
	goose -dir ./migrations postgres "postgresql://test:test@localhost:5432/test?sslmode=disable" down

dt0:
	goose -dir ./migrations postgres "postgresql://test:test@localhost:5432/test?sslmode=disable" down-to 0
