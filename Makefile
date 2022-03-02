.PHONY: create_migration

create_migration:
	migrate create -ext sql -dir cmd/marvin/migrations -seq "$@"

