.PHONY: create_migration tools-lint lint

create_migration:
	migrate create -ext sql -dir cmd/marvin/migrations -seq "$@"

tools-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run