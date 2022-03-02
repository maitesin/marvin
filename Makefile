create_migration:
	migrate create -ext sql -dir devops/db/migrations -seq "$@"
