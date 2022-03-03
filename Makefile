.PHONY: create_migration lint

create_migration:
	migrate create -ext sql -dir cmd/marvin/migrations -seq "$@"

lint:
	golangci-lint run